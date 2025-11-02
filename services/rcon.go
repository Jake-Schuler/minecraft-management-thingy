package services

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorcon/rcon"
)

func GetUsers(conn *rcon.Conn) ([]string, error) {
	response, err := conn.Execute("whitelist list")
	if err != nil {
		return nil, err
	}
	users := parseUsers(response)
	return users, nil
}

func parseUsers(response string) []string {
	// Keep backwards-compatible behaviour: return only the parsed players slice.
	lines := strings.Split(response, "\n")
	w, _ := ParseWhitelist(lines)
	if w.Players == nil {
		return []string{}
	}
	return w.Players
}

// WhitelistResult represents the parsed information from the whitelist output.
type WhitelistResult struct {
	Count   int      `json:"count"`
	Seen    int      `json:"seen"`
	Players []string `json:"players"`
}

// ParseWhitelist parses RCON output like:
// [21:50:28 INFO] There are 3 (out of 5 seen) whitelisted players:
// [21:50:28 INFO] user1, user2, user3, user4 and user 5
// It returns a WhitelistResult with the count, seen and player list.
func ParseWhitelist(lines []string) (WhitelistResult, error) {
	var result WhitelistResult
	// ensure Players is non-nil so JSON marshals to [] instead of null
	result.Players = []string{}

	// Regex to capture "There are <count> (out of <seen> seen)? whitelisted players: <players>"
	// Capturing groups:
	// 1 = count, 2 = seen (optional), 3 = player list (optional)
	// Use a case-insensitive match by prefixing with (?i) so variations in case are handled.
	// Accept "player", "players", or "player(s)" (some servers print the '(s)')
	reCountAndPlayers := regexp.MustCompile(`(?i)There are\s+(\d+)(?:\s*\(out of\s*(\d+)\s*seen\))?\s+whitelisted\s+player(?:\(s\))?s?:\s*(.*)`)

	var foundCount bool
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// The previous regex had a problem because it was expecting the full
		// "[22:28:51 INFO]:" prefix.
		// Let's remove any common prefixes before applying our more focused regex.
		// This makes the regex itself simpler and more focused on the core data.
		// We'll strip common log prefixes.
		prefixIndex := strings.Index(trimmedLine, "] ")
		if prefixIndex != -1 {
			// Check if the character after `]` is a colon, which is also common
			if len(trimmedLine) > prefixIndex+2 && trimmedLine[prefixIndex+1] == ':' {
				trimmedLine = trimmedLine[prefixIndex+3:] // Skip "] :"
			} else {
				trimmedLine = trimmedLine[prefixIndex+2:] // Skip "] "
			}
		}

		matches := reCountAndPlayers.FindStringSubmatch(trimmedLine)
		if matches != nil {
			// matches[1] = count
			if len(matches) >= 2 && matches[1] != "" {
				if count, err := strconv.Atoi(matches[1]); err == nil {
					result.Count = count
				}
			}
			// matches[2] = seen (optional)
			if len(matches) >= 3 && matches[2] != "" {
				if seen, err := strconv.Atoi(matches[2]); err == nil {
					result.Seen = seen
				}
			}
			// matches[3] = player list (optional)
			var playerString string
			if len(matches) >= 4 {
				playerString = strings.TrimSpace(matches[3])
			}
			if playerString != "" {
				// Normalize " and " (before the last name) into commas, then split by comma
				playerString = strings.ReplaceAll(playerString, " and ", ", ")
				rawPlayers := strings.Split(playerString, ",")
				for _, p := range rawPlayers {
					name := strings.TrimSpace(p)
					if name != "" {
						result.Players = append(result.Players, name)
					}
				}
				// We found players on the same line as the count — done.
				break
			}

			// Count line found but no players on the same line. Continue to next lines
			// and look for the players on the following line(s).
			foundCount = true
			continue
		}

		// If we previously saw the count line and this line looks like a list of players,
		// parse it. We identify such a line by ensuring it doesn't contain metadata words
		// like "there are", "whitelist", or "seen" and contains either commas or words.
		if foundCount {
			lower := strings.ToLower(trimmedLine)
			if trimmedLine != "" && !strings.Contains(lower, "there are") && !strings.Contains(lower, "whitelist") && !strings.Contains(lower, "whitelisted player") && !strings.Contains(lower, "seen") {
				// Normalize " and " to commas then split
				linePlayers := strings.ReplaceAll(trimmedLine, " and ", ", ")
				rawPlayers := strings.Split(linePlayers, ",")
				for _, p := range rawPlayers {
					name := strings.TrimSpace(p)
					if name != "" {
						result.Players = append(result.Players, name)
					}
				}
				break
			}
		}
	}

	return result, nil
}

// WhitelistToJSON marshals a WhitelistResult to JSON.
func WhitelistToJSON(w WhitelistResult) ([]byte, error) {
	return json.Marshal(w)
}

func ResetUserPassword(conn *rcon.Conn, username, password string) error {
	cmd := "authme password " + username + " " + password
	_, err := conn.Execute(cmd)
	return err
}

func KickAllPlayers(conn *rcon.Conn) error {
	players, _ := GetUsers(conn)
	for _, player := range players {
		_, err := conn.Execute("minecraft:kick " + player)
		if err != nil {
			return err
		}
	}
	return nil
}

func KickPlayer(conn *rcon.Conn, username string) error {
	_, err := conn.Execute("minecraft:kick " + username)
	return err
}
