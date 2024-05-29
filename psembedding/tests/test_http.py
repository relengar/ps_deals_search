from fastapi.testclient import TestClient
from psembedding.http.healthz import Dependeneies, create_http_server
from psembedding.nats.nats import NatsClient

app = create_http_server(Dependeneies(nats=NatsClient()))
client = TestClient(app)

def test_healthz():
    response = client.get("/healthz")
    assert response.status_code == 200
