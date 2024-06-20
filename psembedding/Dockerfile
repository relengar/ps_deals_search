# syntax=docker/dockerfile:1
ARG PYTHON_VERSION=3.10.12
FROM python:${PYTHON_VERSION}-slim as base

# Prevents Python from writing pyc files.
ENV PYTHONDONTWRITEBYTECODE=1

# Keeps Python from buffering stdout and stderr to avoid situations where
# the application crashes without emitting any logs due to buffering.
ENV PYTHONUNBUFFERED=1

WORKDIR /app

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    embedder

COPY ./pyproject.toml .
RUN --mount=type=cache,target=/root/.cache/pip \
    python -m pip install hatch &&\
    hatch dep show requirements > ./requirements.txt &&\
    python -m pip install -r requirements.txt

RUN mkdir .cache && chown embedder .cache & mkdir .tcache && chown embedder .tcache

USER embedder

COPY ./src .

EXPOSE ${PORT}

CMD ["python", "-m", "psembedding"]
