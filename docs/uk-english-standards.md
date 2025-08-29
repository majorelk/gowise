# UK English Standards for GoWise

GoWise enforces UK English spelling throughout all documentation, comments, and user-facing text via automated CI checks.

## Checked Spellings

The following US spellings are automatically detected and should be replaced with UK equivalents:

| US Spelling | UK Spelling | US Spelling | UK Spelling |
|-------------|-------------|-------------|-------------|
| behavior    | behaviour   | initialize  | initialise  |
| organize    | organise    | prioritize  | prioritise  |
| optimize    | optimise    | analyze     | analyse     |
| catalog     | catalogue   | license (US)| licence (UK)|
| color       | colour      | honor       | honour      |
| labor       | labour      | favor       | favour      |
| traveled    | travelled   | canceled    | cancelled   |
| modeling    | modelling   | leveling    | levelling   |

## Exemptions

The spelling checker intelligently excludes:

### Standard Library Compatibility
```go
// These are allowed - required for Go stdlib
ctx, cancel := context.WithCancel(ctx)
req.Header.Set("Authorization", token)
```

### API Compatibility  
```go
// JSON/XML struct tags may use US spellings for external APIs
type Info struct {
    Licence    string `json:"license"`    // OK - external API requirement
    LicenceUrl string `json:"licenseUrl"` // OK - external API requirement  
}
```

### Import Statements
```go
import "context" // OK - standard library package name
```

## Running the Check

### Locally
```bash
./scripts/check-uk-spelling.sh
```

### In CI
The check runs automatically on every PR and push via GitHub Actions.

## Philosophy

The spelling checker focuses on **comments and documentation** where developers have full control over language choices, while being practical about:

- **Standard library integration** - required US spellings in Go's ecosystem
- **External API compatibility** - JSON/XML specifications that mandate US spellings
- **Existing code tolerance** - only flags new violations to avoid breaking working code

This ensures GoWise maintains consistent UK English presentation while remaining compatible with the broader Go ecosystem.