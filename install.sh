#!/bin/sh

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Configuration ---
APP_NAME="shepherd"
REPO_URL="https://github.com/renoinn/shouldupdate.git"
TMP_DIR="/tmp/${APP_NAME}_install_$$" # Temporary directory for cloning

# --- Helper Functions ---
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

print_success() {
    printf "[32m%s[0m
" "$1"
}

print_warning() {
    printf "[33m%s[0m
" "$1"
}

print_error() {
    printf "[31m%s[0m
" "$1" >&2
}

# --- Pre-flight Checks ---
echo "Checking dependencies..."
if ! command_exists git; then
    print_error "Error: git is not installed. Please install git and try again."
    exit 1
fi

if ! command_exists go; then
    print_error "Error: go is not installed. Please install Go and try again."
    exit 1
fi
echo "Dependencies met."

# --- Determine Installation Directory ---
INSTALL_DIR=""
MAN_DIR=""

# Try /usr/local/bin first
if [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
    MAN_DIR="/usr/local/share/man/man1"
    echo "Will attempt to install to /usr/local/bin"
else
    echo "/usr/local/bin not writable or does not exist."
    # Fallback to ~/.local/bin
    LOCAL_BIN_DIR="${HOME}/.local/bin"
    LOCAL_MAN_DIR="${HOME}/.local/share/man/man1"

    if [ ! -d "${LOCAL_BIN_DIR}" ]; then
        echo "Creating directory: ${LOCAL_BIN_DIR}"
        mkdir -p "${LOCAL_BIN_DIR}"
    fi
    if [ ! -w "${LOCAL_BIN_DIR}" ]; then
        print_error "Error: Cannot write to ${LOCAL_BIN_DIR}. Please check permissions."
        exit 1
    fi
    INSTALL_DIR="${LOCAL_BIN_DIR}"
    MAN_DIR="${LOCAL_MAN_DIR}"
    echo "Will install to ${INSTALL_DIR}"
fi

# --- Installation ---
echo "Starting installation of ${APP_NAME}..."

# 1. Clone repository
echo "Cloning repository from ${REPO_URL}..."
rm -rf "${TMP_DIR}" # Remove temp dir if it exists from a previous failed attempt
git clone --depth 1 "${REPO_URL}" "${TMP_DIR}"
cd "${TMP_DIR}"
echo "Repository cloned."

# 2. Build application
echo "Building application..."
# Assuming main.go is at the root of the repository
# If your main.go is in a subdirectory, adjust the path e.g., cmd/app/main.go
go build -o "${APP_NAME}" main.go
if [ ! -f "${APP_NAME}" ]; then
    print_error "Error: Build failed. ${APP_NAME} binary not found."
    exit 1
fi
echo "Build successful."

# 3. Grant execute permission
echo "Setting execute permission for ${APP_NAME}..."
chmod +x "${APP_NAME}"
echo "Execute permission set."

# 4. Move binary to installation directory
echo "Installing ${APP_NAME} to ${INSTALL_DIR}..."
mv "${APP_NAME}" "${INSTALL_DIR}/"
print_success "${APP_NAME} installed successfully to ${INSTALL_DIR}/${APP_NAME}"

# 5. Install man page (if exists)
# This script looks for a man page at "man/shouldupdate.1" in the cloned repository.
# If your man page is named differently or located elsewhere, adjust MAN_PAGE_SRC.
MAN_PAGE_SRC="${TMP_DIR}/man/${APP_NAME}.1"
if [ -f "${MAN_PAGE_SRC}" ]; then
    echo "Found man page: ${MAN_PAGE_SRC}"
    if [ ! -d "${MAN_DIR}" ]; {
        echo "Creating man directory: ${MAN_DIR}"
        # Create parent directories as needed and don't error if it already exists.
        mkdir -p "${MAN_DIR}"
    }
    fi
    if [ -w "${MAN_DIR}" ]; then
        echo "Installing man page to ${MAN_DIR}..."
        cp "${MAN_PAGE_SRC}" "${MAN_DIR}/"
        print_success "Man page installed to ${MAN_DIR}/${APP_NAME}.1"
    else
        print_warning "Warning: Man page directory ${MAN_DIR} is not writable. Skipping man page installation."
    fi
else
    echo "No man page found at ${MAN_PAGE_SRC} (looked in the 'man' directory of the repository). Skipping man page installation."
fi

# --- Post-installation ---
# 6. PATH notification
if ! echo ":${PATH}:" | grep -q ":${INSTALL_DIR}:"; then
    print_warning "
--- IMPORTANT ---"
    print_warning "The directory ${INSTALL_DIR} is not in your PATH."
    echo "To run '${APP_NAME}' directly, you need to add it to your PATH."
    echo "You can do this by adding one of the following lines to your shell configuration file (e.g., ~/.bashrc, ~/.zshrc, or ~/.profile):"
    echo ""
    echo "  export PATH="${INSTALL_DIR}:\$PATH""
    echo ""
    echo "After adding the line, restart your shell or source the configuration file, e.g.:"
    echo "  source ~/.bashrc  # For bash users"
    echo "  source ~/.zshrc  # For zsh users"
else
    print_success "
${APP_NAME} is ready to use!"
fi

# --- Cleanup ---
echo "Cleaning up temporary files..."
# Important: Go back to the original directory before removing TMP_DIR
# if the script was sourced or if trap EXIT is used.
# However, with set -e, if cd fails, script exits.
# If git clone fails, script exits.
# If build fails, script exits.
# So, we should be in TMP_DIR.
original_dir=$(pwd)
cd / # Go to a safe directory before removing TMP_DIR
rm -rf "${TMP_DIR}"
cd "${original_dir}" # Go back to original dir (though script is about to exit)
echo "Cleanup complete."

print_success "
Installation of ${APP_NAME} finished!"

exit 0
