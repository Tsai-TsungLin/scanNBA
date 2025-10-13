# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

NBA scanner tool that fetches daily NBA game schedules, injury reports, and betting line performance from various web sources. The tool outputs game information with Chinese team translations and injury status in Chinese.

## Building and Running

```bash
# Install dependencies
go mod download

# Build the binary
go build -o nba-scan

# Run all games for today
./nba-scan

# Run games starting at a specific time
./nba-scan --time 11:00
```

## Architecture

This project is currently in transition from a monolithic structure to a modular architecture:

- **Legacy code**: `nba.go` contains the original implementation with all functions (Schedule struct, TeamInit, getInjury, PKTeam, etc.)
- **New structure** (partially complete):
  - `cmd/root.go`: Cobra CLI command definitions
  - `internal/models/`: Data structures (Schedule, TeamMap)
  - `internal/crawler/`: Web scraping logic (injury data from ESPN)
  - `internal/logic/`: Business logic (game analysis, output formatting)

### Key Components

**Data Sources**:
- NBA game schedules: `https://in.global.nba.com/stats2/scores/daily.json`
- Injury reports: `https://www.espn.com/nba/injuries` (scraped with goquery)
- Betting lines: `https://nba.titan007.com/jsData/letGoal/` (regex parsing of JS files)

**Core Functions** (in legacy `nba.go`):
- `PKTeam()`: Main function to fetch and display all games for today
- `PKTeamOnStartTime(st string)`: Filter games by start time
- `getInjury(searchTeam string)`: Scrape injury report for a specific team
- `getDish(searchTeam string)`: Fetch recent betting performance (过盘状况)
- `TeamInit()`: Returns map of English team names to Chinese translations

**Special Cases**:
- LA Clippers team name requires special handling (sometimes "Los Angeles Clippers", sometimes "LA Clippers")
- Time zones: API uses ET (UTC-4), output converts to local time (+13 hours for display)

## Code Migration Notes

The codebase is being refactored from `nba.go` to the `cmd/` and `internal/` structure. Currently:

- The new `internal/logic/analyzer.go` has a simplified `PKTeam()` that's missing:
  - Betting line analysis (`getDish` functionality)
  - Formatted injury comments for Chinese output
  - Complete time filtering in `PKTeamOnStartTime`

- `nba.go` still contains the complete, working implementation

When making changes, consider whether to update the legacy code or continue the migration to the new structure.

## Dependencies

- `github.com/PuerkitoBio/goquery`: HTML parsing for injury data
- `github.com/spf13/cobra`: CLI framework
- Standard library: `regexp`, `time`, `net/http`, `encoding/json`
