import grpc
import service_pb2
import service_pb2_grpc

_HOST = "127.0.0.1"
_PORT = "41005"

def main():
    with grpc.insecure_channel("{0}:{1}".format(_HOST, _PORT)) as channel:
        client = service_pb2_grpc.SayHelloServiceStub(channel=channel)
        response = client.SayHello(service_pb2.SayHelloRequest(name="Tony Stack"))
    print("received: " + response.result)


if __name__ == '__main__':
    main()