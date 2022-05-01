package pleg

import (
	"fmt"
	"minik8s/apiObject"
	"minik8s/kubelet/src/podutil"
	"minik8s/kubelet/src/runtime/container"
	"minik8s/kubelet/src/runtime/runtime"
	"minik8s/kubelet/src/status"
	"minik8s/kubelet/src/types"
	"time"
)

const (
	eventChannelSize      = 10
	relistIntervalSeconds = 10
)

// PodLifecycleEventType define the event type of pod life cycle events.
type PodLifecycleEventType string

const (
	// ContainerStarted - event type when the new state of container is running.
	ContainerStarted PodLifecycleEventType = "ContainerStarted"
	// ContainerDied - event type when the new state of container is exited.
	ContainerDied PodLifecycleEventType = "ContainerDied"
	// ContainerRemoved - event type when the old state of container is exited.
	ContainerRemoved PodLifecycleEventType = "ContainerRemoved"
	// ContainerNeedStart - event type when the container is needed to start.
	ContainerNeedStart PodLifecycleEventType = "ContainerNeedStart"
	// ContainerNeedRestart - event type when the container needs to restart.
	ContainerNeedRestart PodLifecycleEventType = "ContainerNeedRestart"
	// ContainerNeedCreateAndStart - event type when the container needs to create and start.
	ContainerNeedCreateAndStart PodLifecycleEventType = "ContainerNeedCreateAndStart"
	// ContainerNeedRemove - event type when the container needs to be removed.
	ContainerNeedRemove PodLifecycleEventType = "ContainerNeedRemove"
	// PodSync is used to trigger syncing of a pod when the observed change of
	// the state of the pod cannot be captured by any single event above.
	PodSync PodLifecycleEventType = "PodSync"
	// ContainerChanged - event type when the new state of container is unknown.
	ContainerChanged PodLifecycleEventType = "ContainerChanged"
)

// PodLifecycleEvent is an event that reflects the change of the pod state.
type PodLifecycleEvent struct {
	// The pod ID.
	ID types.UID
	// The api object pod itself Pod
	Pod *apiObject.Pod
	// The type of the event.
	Type PodLifecycleEventType
	// The accompanied data which varies based on the event type.
	//   - ContainerStarted/ContainerStopped: the container name (string).
	//   - All other event types: unused.
	Data interface{}
}

type podStatusRecord struct {
	OldStatus     *runtime.PodStatus
	CurrentStatus *runtime.PodStatus
}

type podStatusRecords map[types.UID]*podStatusRecord

func (statusRecords podStatusRecords) UpdateRecord(podUID types.UID, newStatus *runtime.PodStatus) {
	if record, exists := statusRecords[podUID]; exists {
		record.OldStatus = record.CurrentStatus
		record.CurrentStatus = newStatus
	} else {
		statusRecords[podUID] = &podStatusRecord{
			OldStatus:     nil,
			CurrentStatus: newStatus,
		}
	}
}

func (statusRecords podStatusRecords) RemoveRecord(podUID types.UID) {
	delete(statusRecords, podUID)
}

func (statusRecords podStatusRecords) GetRecord(podUID types.UID) *podStatusRecord {
	return statusRecords[podUID]
}

type PodRestartContainerArgs struct {
	ContainerID       container.ContainerID
	ContainerFullName string
}

type Manager interface {
	Updates() chan *PodLifecycleEvent
	Start()
}

func NewPlegManager(statusManager status.Manager, podManager runtime.Manager) Manager {
	return &manager{
		eventCh:          make(chan *PodLifecycleEvent, eventChannelSize),
		statusManager:    statusManager,
		podManager:       podManager,
		podStatusRecords: make(podStatusRecords),
	}
}

type manager struct {
	eventCh          chan *PodLifecycleEvent
	statusManager    status.Manager
	podManager       runtime.Manager
	podStatusRecords podStatusRecords
}

func newPodLifecycleEvent(podUID types.UID, pod *apiObject.Pod, eventType PodLifecycleEventType, data interface{}) *PodLifecycleEvent {
	return &PodLifecycleEvent{
		ID:   podUID,
		Pod:  pod,
		Type: eventType,
		Data: data,
	}
}

func (m *manager) addStartedLifecycleEvent(podUID types.UID, pod *apiObject.Pod, containerID container.ContainerID) {
	m.eventCh <- newPodLifecycleEvent(podUID, pod, ContainerStarted, containerID)
}

func (m *manager) addNeedRemoveLifecycleEvent(podUID types.UID, pod *apiObject.Pod, containerID container.ContainerID) {
	m.eventCh <- newPodLifecycleEvent(podUID, pod, ContainerNeedRemove, containerID)
}

func (m *manager) addNeedRestartLifecycleEvent(podUID types.UID, pod *apiObject.Pod, args PodRestartContainerArgs) {
	m.eventCh <- newPodLifecycleEvent(podUID, pod, ContainerNeedRestart, args)
}

func (m *manager) addNeedStartLifecycleEvent(podUID types.UID, pod *apiObject.Pod, args PodRestartContainerArgs) {
	m.eventCh <- newPodLifecycleEvent(podUID, pod, ContainerNeedStart, args)
}

func (m *manager) addNeedCreateAndStartLifecycleEvent(podUID types.UID, pod *apiObject.Pod, target *apiObject.Container) {
	m.eventCh <- newPodLifecycleEvent(podUID, pod, ContainerNeedCreateAndStart, target)
}

func (m *manager) addDiedLifecycleEvent(podUID types.UID, pod *apiObject.Pod, containerID container.ContainerID) {
	m.eventCh <- newPodLifecycleEvent(podUID, pod, ContainerDied, containerID)
}

func (m *manager) addRemovedLifecycleEvent(podUID types.UID, pod *apiObject.Pod, containerID container.ContainerID) {
	m.eventCh <- newPodLifecycleEvent(podUID, pod, ContainerRemoved, containerID)
}

func (m *manager) addPodSyncLifecycleEvent(podUID types.UID, pod *apiObject.Pod, containerID container.ContainerID) {
	m.eventCh <- newPodLifecycleEvent(podUID, pod, PodSync, containerID)
}

func (m *manager) addChangedLifecycleEvent(podUID types.UID, pod *apiObject.Pod, containerID container.ContainerID) {
	m.eventCh <- newPodLifecycleEvent(podUID, pod, ContainerChanged, containerID)
}

func (m *manager) removeAllContainers(runtimePodStatus *runtime.PodStatus) {
	for _, cs := range runtimePodStatus.ContainerStatuses {
		m.addNeedRemoveLifecycleEvent(runtimePodStatus.ID, nil, cs.ID)
	}
}

// compareAndProduceLifecycleEvents compares given runtime pod statuses
// with pod api object, and produce corresponding lifecycle events
/// TODO what about pause?
func (m *manager) compareAndProduceLifecycleEvents(apiPod *apiObject.Pod, runtimePodStatus *runtime.PodStatus) {
	podUID := runtimePodStatus.ID
	m.podStatusRecords.UpdateRecord(podUID, runtimePodStatus)
	record := m.podStatusRecords.GetRecord(podUID)
	oldStatus, currentStatus := record.OldStatus, record.CurrentStatus

	// apiPod == nil means the pod is no longer existent, skip this
	if apiPod == nil {
		return
	}

	notIncludedContainerNameMap := make(map[string]struct{})
	for _, c := range apiPod.Containers() {
		notIncludedContainerNameMap[c.Name] = struct{}{}
	}

	for _, cs := range currentStatus.ContainerStatuses {
		parseSucc, containerName, _, _, _, _ := podutil.ParseContainerFullName(cs.Name)
		// illegal containerName, need remove it
		if !parseSucc {
			//m.addNeedRemoveLifecycleEvent(podUID, cs.ID)
			continue
		}

		// Only deal with it when state has changed
		needDealWith := oldStatus == nil
		if !needDealWith {
			oldCs := oldStatus.GetContainerStatusByName(cs.Name)
			needDealWith = oldCs == nil || oldCs.State != cs.State
		}

		if needDealWith {
			switch cs.State {
			case container.ContainerStateRunning:
			case container.ContainerStateCreated:
				//if apiPod.GetContainerByName(containerName) == nil {
				//	m.addNeedRemoveLifecycleEvent(podUID, cs.ID)
				//}
				break
			// Need restart it
			case container.ContainerStateExited:
				if apiPod.GetContainerByName(containerName) != nil {
					m.addNeedRestartLifecycleEvent(podUID, apiPod, PodRestartContainerArgs{cs.ID, cs.Name})
				}
			default:
				m.addChangedLifecycleEvent(podUID, apiPod, cs.ID)
			}
		}
		delete(notIncludedContainerNameMap, containerName)
	}
	// Need to create all the container that has not been created
	for notIncludeContainerName := range notIncludedContainerNameMap {
		m.addNeedCreateAndStartLifecycleEvent(podUID, apiPod, apiPod.GetContainerByName(notIncludeContainerName))
	}
}

func (m *manager) relist() error {
	// Step 1: Get all *runtime* pod statuses
	runtimePodStatuses, err := m.podManager.GetPodStatuses()
	if err != nil {
		return err
	}

	// Step 2: Get pod api object, and according to the api object, produce lifecycle events
	var apiPod *apiObject.Pod
	for podUID, runtimePodStatus := range runtimePodStatuses {
		apiPod = m.statusManager.GetPod(podUID)
		m.compareAndProduceLifecycleEvents(apiPod, runtimePodStatus)
	}

	return nil
}

func (m *manager) run() {
	ticker := time.Tick(relistIntervalSeconds * time.Second)
	for {
		select {
		case <-ticker:
			if err := m.relist(); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func (m *manager) Updates() chan *PodLifecycleEvent {
	return m.eventCh
}

func (m *manager) Start() {
	go m.run()
}
