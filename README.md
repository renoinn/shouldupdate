# shouldupdate

`shouldupdate` is a command-line utility that [TODO: Add a brief description of what shouldupdate does. e.g., 'checks if your system packages need an update'].

## Installation

You can install `shouldupdate` using the provided installation script. This script will clone the repository, build the binary, and install it to a standard location on your system (Linux or macOS).

**Prerequisites:**
*   `git`
*   `go` (Go programming language compiler)

**Steps:**

1.  **Download the installation script:**

    You can download the script using `curl` or `wget`.

    Using `curl`:
    ```bash
    curl -LO https://raw.githubusercontent.com/renoinn/shouldupdate/main/install.sh
    ```
    *(Note: The URL assumes the `install.sh` script will be in the `main` branch of your repository. Adjust if it's located elsewhere, e.g., a `release` branch or a specific tag.)*

    Using `wget`:
    ```bash
    wget https://raw.githubusercontent.com/renoinn/shouldupdate/main/install.sh
    ```
    *(Same note as above regarding the URL)*

2.  **Make the script executable (if necessary):**
    ```bash
    chmod +x install.sh
    ```

3.  **Run the installation script:**
    ```bash
    ./install.sh
    ```
    Or, if you prefer:
    ```bash
    sh install.sh
    ```

    The script will:
    *   Check for `git` and `go`.
    *   Clone the latest version of the `shouldupdate` repository.
    *   Build the `shouldupdate` binary.
    *   Attempt to install the binary to `/usr/local/bin`. If that's not possible (e.g., due to permissions), it will try to install to `~/.local/bin` (and create this directory if it doesn't exist).
    *   If a man page (`man/shouldupdate.1`) is found in the repository, it will attempt to copy it to the corresponding man directory (`/usr/local/share/man/man1` or `~/.local/share/man/man1`).
    *   Notify you if the installation directory is not in your `PATH`.

**Post-Installation: Updating your PATH (if needed)**

If the installation script informs you that the installation directory (e.g., `~/.local/bin`) is not in your `PATH`, you'll need to add it to be able to run `shouldupdate` directly from your terminal.

Add one of the following lines to your shell's configuration file:

*   For **bash** users (commonly `~/.bashrc` or `~/.bash_profile` on macOS):
    ```bash
    export PATH="$HOME/.local/bin:$PATH"
    ```
    (If installed to `/usr/local/bin`, it's usually already in `PATH`. If not, use `export PATH="/usr/local/bin:$PATH"`)

*   For **zsh** users (commonly `~/.zshrc`):
    ```bash
    export PATH="$HOME/.local/bin:$PATH"
    ```
    (If installed to `/usr/local/bin`, it's usually already in `PATH`. If not, use `export PATH="/usr/local/bin:$PATH"`)

After adding the line, either restart your terminal or source the configuration file (e.g., `source ~/.bashrc` or `source ~/.zshrc`).

## Uninstallation

To uninstall `shouldupdate`:

1.  **Remove the binary:**
    *   Identify where `shouldupdate` was installed. The installation script mentions this, or you can use `which shouldupdate`.
    *   Delete the binary. For example:
        ```bash
        # If installed in ~/.local/bin
        rm ~/.local/bin/shouldupdate

        # Or, if installed in /usr/local/bin
        # sudo rm /usr/local/bin/shouldupdate
        ```

2.  **Remove the man page (if installed):**
    *   If a man page was installed, remove it from the man directory.
        ```bash
        # If man page installed in ~/.local/share/man/man1
        # rm ~/.local/share/man/man1/shouldupdate.1

        # Or, if man page installed in /usr/local/share/man/man1
        # sudo rm /usr/local/share/man/man1/shouldupdate.1
        ```

## Usage

[TODO: Add instructions on how to use shouldupdate. Include common commands and examples.]

```bash
shouldupdate [options]
```

---

This `README.md` provides the necessary information for users to install, configure, and uninstall the `shouldupdate` application.
Remember to fill in the `TODO` sections with specific details about your application.
