## Git
Always do git add . and check if there are some things that shouldn't be committed (if so add them to gitignore). than create commits and branch based on all the changes made.

## Git branches
Use Prefixes to Indicate Purpose. Use kebab-case. Branch names should be concise yet informative. A good branch name briefly describes what it is for without being overly long or vague.

- `feat:` - For new features or functionalities.
- `bug:` - For fixing bugs in the code.
- `hotfix:` - For urgent patches, usually applied to production.
- `refactor:` - For improving code structure without changing functionality.
- `test:` - For writing or improving automated tests.
- `doc:` - For documentation updates.

Examples: 

    feature/user-authentication
    bugfix/fix-login-error
    hotfix/urgent-patch-crash
    design/update-navbar
    refactor/remove-unused-code
    test/add-unit-tests
    doc/update-readme



## Git Commits
Use conventional commit prefixes for all commit messages:

- `feat:` - New features or functionality
- `fix:` - Bug fixes
- `refactor:` - Code refactoring without changing functionality
- `chore:` - Maintenance tasks, dependency updates, config changes
- `cicd:` - CI/CD pipeline changes

Example: `feat: add dark mode toggle`

## Releasing (GoReleaser + Homebrew)

To publish a new release and update the Homebrew tap (`Tiimie1/homebrew-tap`):

1. Tag the release (GoReleaser derives the version from git tags):
   ```bash
   git tag -a v0.x.x -m "Release v0.x.x"
   git push origin v0.x.x
   ```

2. Run GoReleaser (requires a `GITHUB_TOKEN` with repo write access):
   ```bash
   GITHUB_TOKEN=<your-token> goreleaser release --clean
   ```
   If the token is already exported in your shell, just run:
   ```bash
   goreleaser release --clean
   ```

GoReleaser will build binaries for darwin amd64/arm64, create a GitHub release, and push the updated formula to the Homebrew tap automatically.
