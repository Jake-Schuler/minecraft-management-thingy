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

	// Regex to capture "There are <count> whitelisted player(s): <players>"
	// This regex is more flexible, ignoring the timestamp/INFO prefix and focusing
	// on "There are X whitelisted player(s):" and then capturing the player list.
	reCountAndPlayers := regexp.MustCompile(`There are\s+(\d+)\s+whitelisted player\(s\):\s*(.*)`)

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

		if matches := reCountAndPlayers.FindStringSubmatch(trimmedLine); matches != nil {
			if count, err := strconv.Atoi(matches[1]); err == nil {
				result.Count = count
			}

			playerString := strings.TrimSpace(matches[2])
			if playerString != "" {
				// Split by comma and trim each player name
				rawPlayers := strings.Split(playerString, ",")
				for _, p := range rawPlayers {
					name := strings.TrimSpace(p)
					if name != "" {
						result.Players = append(result.Players, name)
					}
				}
			}
			// Assuming only one line contains this information in this format
			// If your output can have multiple such lines, remove this break.
			break
		}
	}

	return result, nil
}

// WhitelistToJSON marshals a WhitelistResult to JSON.
func WhitelistToJSON(w WhitelistResult) ([]byte, error) {
	return json.Marshal(w)
}

func ResetUserPassword(conn *rcon.Conn, username, password string) error {
	cmd := "say authme password " + username + " " + password
	_, err := conn.Execute(cmd)
	return err
}
