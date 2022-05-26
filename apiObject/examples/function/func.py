def main(params):
    name = params["name"]
    a = params["a"]
    b = params["b"]

    return "Hello world, "+name+","+str(a)+"+"+str(b)+"="+str(a+b)+"\n"
