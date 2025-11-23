FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

# Install Ollama
RUN curl -fsSL https://ollama.com/install.sh | sh

# Expose Ollama port
EXPOSE 11434

# We don't start Ollama here because the entrypoint of the devcontainer 
# usually handles the shell. We'll handle starting Ollama in postCreateCommand or a script.
