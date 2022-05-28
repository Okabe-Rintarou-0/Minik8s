package kpa

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"minik8s/apiObject"
	"minik8s/entity"
)

type task struct {
	Function string
	Next     *node
}

type nodeType byte

const (
	taskType nodeType = iota
	choiceType
)

type judgeFunc func(data entity.FunctionData) bool

type branch struct {
	Variable  string
	JudgeFunc judgeFunc
	Next      *node
}

func (b *branch) Satisfied(data entity.FunctionData) bool {
	return b.JudgeFunc(data)
}

type choice struct {
	Branches []branch
}

func (c *choice) ChooseSatisfied(data entity.FunctionData) *node {
	for i, br := range c.Branches {
		if br.Satisfied(data) {
			return c.Branches[i].Next
		}
	}
	return nil
}

type node struct {
	Type nodeType
	*task
	*choice
}

type dag struct {
	Root        *node
	EntryParams string
}

func toJudgeFunc(c *apiObject.Choice) judgeFunc {
	variable := c.Variable

	switch {
	case c.NumericEquals != nil:
		{
			return func(data entity.FunctionData) bool {
				res := gjson.Get(string(data), variable)
				return res.Int() == *c.NumericNotEquals
			}
		}
	case c.NumericNotEquals != nil:
		{
			return func(data entity.FunctionData) bool {
				res := gjson.Get(string(data), variable)
				return res.Int() != *c.NumericNotEquals
			}
		}
	}

	return func(data entity.FunctionData) bool {
		return false
	}
}

func Workflow2DAG(wf *apiObject.Workflow) *dag {
	nodeMap := make(map[string]apiObject.WorkflowNode)
	dagMap := make(map[string]*node)
	for nodeName := range wf.Nodes {
		nodeMap[nodeName] = wf.Nodes[nodeName]
		dagMap[nodeName] = &node{}
	}

	if _, exists := nodeMap[wf.StartAt]; !exists {
		return nil
	}

	params, _ := json.Marshal(wf.Params)

	root := buildDAG(wf.StartAt, dagMap, nodeMap)
	if root == nil {
		return nil
	}
	return &dag{
		Root:        root,
		EntryParams: string(params),
	}
}

func buildDAG(curNode string, dagMap map[string]*node, nodeMap map[string]apiObject.WorkflowNode) *node {
	wfNode := nodeMap[curNode]
	dagNode := dagMap[curNode]
	if dagNode == nil {
		return nil
	}

	var next *node = nil
	switch wfNode.Type {
	case apiObject.TaskNode:
		if wfNode.Task != nil && wfNode.Next != nil {
			next = buildDAG(*wfNode.Next, dagMap, nodeMap)
		} else {
			next = nil
		}
		return &node{
			Type: taskType,
			task: &task{
				Function: curNode,
				Next:     next,
			},
			choice: nil,
		}
	case apiObject.ChoiceNode:
		choices := wfNode.Choices

		var branches []branch
		if choices != nil {
			for _, c := range choices.Choices {
				if c.Next != nil {
					next = buildDAG(*c.Next, dagMap, nodeMap)
				} else {
					next = nil
				}
				br := branch{
					Variable:  c.Variable,
					JudgeFunc: toJudgeFunc(&c),
					Next:      next,
				}
				branches = append(branches, br)
			}
		}
		return &node{
			Type: choiceType,
			task: nil,
			choice: &choice{
				Branches: branches,
			},
		}
	}
	return nil
}

func traverse(cur *node) {
	if cur == nil {
		return
	}
	switch cur.Type {
	case taskType:
		fmt.Printf("This is a task node: %s\n", cur.Function)
		traverse(cur.Next)
	case choiceType:
		fmt.Printf("This is a choice node: \n")
		for i, br := range cur.Branches {
			fmt.Printf("If enter branch %d: \n", i)
			traverse(br.Next)
		}
	}
}

func TraverseDAG(dag *dag) {
	traverse(dag.Root)
}
