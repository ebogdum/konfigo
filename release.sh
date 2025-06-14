#!/bin/bash

# --- Configuration ---
# build_dir: Path to the directory containing your compiled binaries/artifacts
# default_branch: Your main development branch (e.g., main, master)
# remote_name: The name of your remote repository (usually "origin")
# prerelease_suffix: Suffix for pre-releases (e.g., "beta", "rc.1"). Leave empty for stable.
# generate_notes: Set to "true" to automatically generate release notes from commits.
#                 Set to "false" to be prompted for notes or use a notes file.
# notes_file: Path to a file containing release notes (if generate_notes is false and you want to use a file).

build_dir="builds" # IMPORTANT: Change this to your actual build output directory
default_branch="main"
remote_name="origin"
prerelease_suffix="" # e.g., "rc.1" or "beta"
generate_notes="false"
notes_file="" # e.g., "RELEASE_NOTES.md"

# --- Helper Functions ---
get_latest_tag() {
    git fetch --tags "$remote_name" >/dev/null 2>&1
    # Get the latest semantic version tag (vX.Y.Z or X.Y.Z)
    git tag --sort=-v:refname | grep -E '^v?[0-9]+\.[0-9]+\.[0-9]+$' | head -n1
}

increment_version() {
    local version="$1"
    local increment_type="$2"
    local current_major current_minor current_patch

    # Remove 'v' prefix if present
    version="${version#v}"

    IFS='.' read -r current_major current_minor current_patch <<< "$version"

    case "$increment_type" in
        patch)
            current_patch=$((current_patch + 1))
            ;;
        minor)
            current_minor=$((current_minor + 1))
            current_patch=0
            ;;
        major)
            current_major=$((current_major + 1))
            current_minor=0
            current_patch=0
            ;;
        *)
            echo "Error: Invalid increment type '$increment_type'. Use 'patch', 'minor', or 'major'."
            exit 1
            ;;
    esac
    echo "v${current_major}.${current_minor}.${current_patch}"
}

# --- Main Script ---

# 0. Prerequisites check
if ! command -v gh &> /dev/null; then
    echo "GitHub CLI 'gh' could not be found. Please install it."
    exit 1
fi
if ! command -v git &> /dev/null; then
    echo "Git 'git' could not be found. Please install it."
    exit 1
fi
if [ ! -d "$build_dir" ]; then
    echo "Error: Build directory '$build_dir' not found."
    echo "Please ensure your project is built and artifacts are in this directory."
    exit 1
fi
if [ -z "$(ls -A "$build_dir")" ]; then
    echo "Error: Build directory '$build_dir' is empty."
    exit 1
fi


# 1. Ensure we are on the default branch and it's clean
echo "Checking Git status..."
git checkout "$default_branch"
if ! git diff --quiet HEAD || ! git diff --cached --quiet HEAD; then
    echo "Your working directory or staging area is not clean. Please commit or stash changes."
    exit 1
fi
echo "Pulling latest changes from $remote_name/$default_branch..."
git pull "$remote_name" "$default_branch"

# 2. Get the latest tag and determine the next version
latest_tag=$(get_latest_tag)
new_version=""

if [ -z "$latest_tag" ]; then
    echo "No existing semantic version tags found."
    read -r -p "Enter the initial version (e.g., v0.1.0 or 0.1.0): " new_version
    if [[ ! "$new_version" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Error: Invalid version format. Please use (v)X.Y.Z (e.g., v0.1.0)."
        exit 1
    fi
    # Ensure 'v' prefix
    if [[ ! "$new_version" =~ ^v ]]; then
        new_version="v${new_version}"
    fi
else
    echo "Latest tag found: $latest_tag"
    read -r -p "Increment type (patch, minor, major) or specify full version (e.g., v1.2.3): " version_input

    if [[ "$version_input" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        new_version="$version_input"
        if [[ ! "$new_version" =~ ^v ]]; then
            new_version="v${new_version}"
        fi
    elif [[ "$version_input" == "patch" || "$version_input" == "minor" || "$version_input" == "major" ]]; then
        new_version=$(increment_version "$latest_tag" "$version_input")
    else
        echo "Error: Invalid input. Use 'patch', 'minor', 'major', or a full version like 'vX.Y.Z'."
        exit 1
    fi
fi

if [ -n "$prerelease_suffix" ]; then
    new_version="${new_version}-${prerelease_suffix}"
fi

echo "New version will be: $new_version"
read -r -p "Proceed with creating tag and release? (y/N): " confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo "Aborted."
    exit 0
fi

# 3. Create and push Git tag
echo "Creating Git tag $new_version..."
if git tag "$new_version"; then
    echo "Pushing tag $new_version to $remote_name..."
    git push "$remote_name" "$new_version"
else
    echo "Error: Failed to create Git tag. Does it already exist?"
    exit 1
fi

# 4. Create GitHub Release and upload artifacts
echo "Creating GitHub Release for $new_version..."
release_title="Release $new_version"
release_options=()

if [ -n "$prerelease_suffix" ]; then
    release_options+=("--prerelease")
    release_title="${release_title} (${prerelease_suffix})"
fi

if [ "$generate_notes" == "true" ]; then
    release_options+=("--generate-notes")
else
    if [ -n "$notes_file" ] && [ -f "$notes_file" ]; then
        release_options+=("--notes-file" "$notes_file")
    else
        read -r -e -p "Enter release notes (or leave blank to open editor): " custom_notes
        if [ -n "$custom_notes" ]; then
            release_options+=("--notes" "$custom_notes")
        fi
    fi
fi

# Construct the list of files to upload
# Using find to handle spaces in filenames and list them properly for gh
files_to_upload=()
while IFS= read -r -d $'\0' file; do
    files_to_upload+=("$file")
done < <(find "$build_dir" -type f -print0)

if [ ${#files_to_upload[@]} -eq 0 ]; then
    echo "Warning: No files found in $build_dir to upload."
else
    echo "Found files to upload:"
    printf " - %s\n" "${files_to_upload[@]}"
fi


if gh release create "$new_version" "${files_to_upload[@]}" --title "$release_title" "${release_options[@]}"; then
    echo "Successfully created GitHub Release $new_version and uploaded artifacts."
    gh release view "$new_version" --web # Open in browser
else
    echo "Error: Failed to create GitHub Release."
    echo "You might need to: "
    echo "  1. Ensure 'gh' is authenticated with sufficient permissions (repo scope)."
    echo "  2. Check if a release for tag '$new_version' already exists."
    echo "The tag '$new_version' was pushed. You might need to create the release manually or delete the tag and retry."
    exit 1
fi

echo "Done."