from concurrent import futures

import grpc

from detect.v1 import detect_pb2
from detect.v1 import detect_pb2_grpc


class FoodDetectService(detect_pb2_grpc.FoodDetectServiceServicer):
    def AnalyzeFood(self, request, context):
        object_key = request.object_key.strip()

        # Mock response for initial integration testing.
        return detect_pb2.AnalyzeFoodResponse(
            items=[
                detect_pb2.FoodItem(
                    name="com tam",
                    confidence=0.93,
                    calories_est=650.0,
                ),
                detect_pb2.FoodItem(
                    name="trung op la",
                    confidence=0.88,
                    calories_est=120.0,
                ),
            ],
            note=f"analyzed from object_key={object_key}",
        )


def serve(host="[::]", port=50051):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    detect_pb2_grpc.add_FoodDetectServiceServicer_to_server(FoodDetectService(), server)
    server.add_insecure_port(f"{host}:{port}")
    server.start()
    print(f"Detect gRPC server listening on {host}:{port}")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
