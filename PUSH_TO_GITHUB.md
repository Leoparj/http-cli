# Push to GitHub - Quick Guide

## Option 1: Automated (if you have GitHub token)

1. Set your GitHub token (get one from https://github.com/settings/tokens):
   ```powershell
   $env:GITHUB_TOKEN = "your_token_here"
   ```

2. Run the script:
   ```powershell
   .\create-repo.ps1
   ```

## Option 2: Manual (Recommended)

### Step 1: Create Repository on GitHub

1. Go to **https://github.com/new**
2. **Repository name**: `http-cli`
3. **Description**: `A colorful, lightweight command-line HTTP client similar to Postman. Built with Go for fast performance and beautiful terminal output.`
4. Choose **Public**
5. **DO NOT** check any initialization options (README, .gitignore, license)
6. Click **"Create repository"**

### Step 2: Push Your Code

Replace `YOUR_USERNAME` with your GitHub username:

```bash
git remote add origin https://github.com/YOUR_USERNAME/http-cli.git
git push -u origin main
```

### Step 3: Add Topics/Tags

After pushing, go to your repository and click the gear icon ⚙️ next to "About", then add these topics:

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

## Repository Information

- **Name**: `http-cli`
- **Description**: A colorful, lightweight command-line HTTP client similar to Postman. Built with Go for fast performance and beautiful terminal output.
- **Language**: Go
- **License**: MIT (you can add a LICENSE file)

## Done! 🎉

Your repository will be available at: `https://github.com/YOUR_USERNAME/http-cli`

