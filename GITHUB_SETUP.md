# GitHub Repository Setup Instructions

## Repository Details

- **Name**: `http-cli`
- **Description**: A colorful, lightweight command-line HTTP client similar to Postman. Built with Go for fast performance and beautiful terminal output.
- **Topics/Tags**: `go`, `cli`, `http-client`, `rest-api`, `postman-alternative`, `terminal-tool`, `http-requests`, `api-testing`, `command-line`, `developer-tools`

## Steps to Push to GitHub

### 1. Create the Repository on GitHub

1. Go to [https://github.com/new](https://github.com/new)
2. **Repository name**: `http-cli`
3. **Description**: `A colorful, lightweight command-line HTTP client similar to Postman. Built with Go for fast performance and beautiful terminal output.`
4. Choose **Public** (or Private if you prefer)
5. **DO NOT** check "Add a README file", "Add .gitignore", or "Choose a license" (we already have these)
6. Click **"Create repository"**

### 2. Add Remote and Push

After creating the repository, run these commands (replace `YOUR_USERNAME` with your GitHub username):

```bash
git remote add origin https://github.com/YOUR_USERNAME/http-cli.git
git branch -M main
git push -u origin main
```

### 3. Add Topics/Tags

After pushing, go to your repository on GitHub and:

1. Click on the gear icon ⚙️ next to "About" section
2. Add these topics (one per line or comma-separated):
   - `go`
   - `cli`
   - `http-client`
   - `rest-api`
   - `postman-alternative`
   - `terminal-tool`
   - `http-requests`
   - `api-testing`
   - `command-line`
   - `developer-tools`

### 4. Optional: Add GitHub Actions Badge

You can add a badge to your README.md:

```markdown
![GitHub release](https://img.shields.io/github/v/release/YOUR_USERNAME/http-cli)
![Go Version](https://img.shields.io/badge/go-1.21+-blue)
![License](https://img.shields.io/badge/license-MIT-green)
```

## Quick Command Reference

```bash
# Add remote (replace YOUR_USERNAME)
git remote add origin https://github.com/YOUR_USERNAME/http-cli.git

# Rename branch to main
git branch -M main

# Push to GitHub
git push -u origin main
```

## Repository Settings to Consider

- Enable GitHub Discussions (for Q&A)
- Add topics/tags as listed above
- Consider adding a LICENSE file (MIT recommended)
- Enable GitHub Actions for CI/CD (optional)

