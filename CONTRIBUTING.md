# Contributing to Cloudflared GUI

Thank you for your interest in contributing to Cloudflared GUI! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and constructive in all interactions. We aim to maintain a welcoming and inclusive community.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in GitHub Issues
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - System information (OS, versions, etc.)
   - Relevant logs or screenshots

### Suggesting Features

1. Check if the feature has already been suggested
2. Create a new issue describing:
   - The problem you're trying to solve
   - Your proposed solution
   - Any alternative solutions considered
   - Additional context

### Pull Requests

1. **Fork the repository** and create a new branch from `main`
2. **Make your changes**:
   - Follow the existing code style
   - Add tests if applicable
   - Update documentation as needed
3. **Test your changes** thoroughly
4. **Commit with clear messages**:
   ```
   feat: add new feature
   fix: resolve bug in component
   docs: update README
   ```
5. **Push to your fork** and submit a pull request
6. **Wait for review** and address any feedback

## Development Setup

See [README.md](README.md) for setup instructions.

## Code Style

### Go Backend

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `go vet` before committing
- Add comments for exported functions

### TypeScript/React Frontend

- Use TypeScript for type safety
- Follow React best practices
- Use functional components with hooks
- Keep components focused and reusable

### CSS

- Use semantic class names
- Support dark mode where applicable
- Ensure mobile responsiveness

## Project Structure

```
cloudflared-gui/
├── apps/
│   ├── backend/         # Go API server
│   └── dashboard/       # React frontend
└── packages/
    ├── types/           # Shared types
    └── ui/              # Shared components
```

## Testing

### Backend

```bash
cd apps/backend
go test ./...
```

### Frontend

```bash
cd apps/dashboard
npm run build  # Ensure it builds without errors
```

## Commit Message Guidelines

Follow conventional commits:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

## Questions?

Feel free to open an issue for any questions or concerns.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

