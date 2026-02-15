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
