def main(params):
    x = params["x"]
    x = x + 5
    result = {
        "x": x
    }
    return json.dumps(result)
