# Camp 🏕️

An open source email marketing platform - a simple, self-hosted alternative to Mailchimp.

## About

Camp is a lightweight email marketing solution built with Go. Send campaigns, manage contacts, and track engagement without the complexity or cost of traditional email marketing platforms.

## Features

- 📧 **Campaign Management** - Create and send email campaigns to your contacts
- 👥 **Contact Management** - Import, organize, and segment your mailing lists
- 📊 **Analytics** - Track opens, clicks, and engagement metrics
- 🎨 **Template System** - Design and reuse email templates
- 🚀 **Self-Hosted** - Full control over your data and infrastructure
- ⚡ **Lightweight** - Built with Go's standard library for minimal dependencies

## Tech Stack

- **Language**: Go
- **HTTP Server**: Go's built-in `net/http` module
- **Database**: SQLite

## Installation

### Prerequisites

- Go 1.25 or higher
- SQLite 3

### Setup

1. Clone the repository:
```bash
git clone https://github.com/codersgyan/camp.git
cd camp
```

2. Install dependencies:
```bash
go mod download
```

3. Configure your environment:
```bash
cp .env.example .env
# Edit .env with your settings
```

4. Run the application:
```bash
// make sure to install `make` before this.
make run
```

The server will start on `http://localhost:8080` (or your configured port).

## Usage

### API Endpoints

(Document your main API endpoints here)

```
POST   /api/campaigns          - Create a new campaign
GET    /api/campaigns/:id      - Get campaign details
POST   /api/contacts           - Add a contact
GET    /api/contacts           - List contacts
POST   /api/contacts/tag       - Add a tags to contact
PATCH  /api/contacts/:id/tag   - Remove a tag from contact
```

## Development

```bash
# Run tests
go test ./...

# Build for production
make build
```

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] Automated email sequences
- [ ] A/B testing
- [ ] Advanced segmentation
- [ ] REST API documentation
- [ ] Web UI dashboard
- [ ] Webhook support

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/codersgyan/camp/issues).

---

Made with 💛 by the open source community