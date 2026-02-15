# Image Slimmer Core

Image Slimmer Core is a deterministic container image analysis engine written in Go, designed to model, interpret, and structure container image metadata into reproducible slimming plans. It acts as the foundational intelligence layer for systems that need deep visibility into image composition without coupling to Docker CLI execution, runtime mutation, or orchestration logic.

The project focuses on analytical correctness and deterministic modeling. Given a container image reference, the engine resolves, normalizes, and interprets its structural data to produce a consistent representation that can be consumed by higher-level systems such as CLIs, APIs, CI/CD pipelines, or governance platforms.

It is not a wrapper around Docker commands. It is not a runtime tracer. It is not an image rebuilder.
It is a computation engine whose sole responsibility is to transform container image metadata into structured, machine-consumable slimming intelligence with predictable, reproducible output.

<img width="300" height="300" alt="ChatGPT Image Feb 14, 2026, 08_56_32 PM" src="https://github.com/user-attachments/assets/7db2e338-6857-4208-b41c-8b6b76bedc08" />
