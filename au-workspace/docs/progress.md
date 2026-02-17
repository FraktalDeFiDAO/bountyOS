# Bounty Progress Tracking

Use `./bin/au bounty ...` to track multiple bounties at once.

## Data model

Tracker file: `config/bounties.tsv`

Fields:
1. `id`
2. `issue_url`
3. `title`
4. `project`
5. `status`
6. `priority`
7. `owner`
8. `progress` (0-100)
9. `next_action`
10. `blocker`
11. `created_at` (UTC)
12. `updated_at` (UTC)

## Allowed values

- `status`: `backlog`, `ready`, `in_progress`, `blocked`, `review`, `submitted`, `paid`, `done`, `dropped`
- `priority`: `low`, `medium`, `high`, `critical`

## Commands

```bash
# list tracked bounties
./bin/au bounty list

# add a bounty
./bin/au bounty add <id> <project> <issue_url> --title "<title>" \
  [--status backlog] [--priority medium] [--owner unassigned] \
  [--progress 0] [--next "-"] [--blocker "-"]

# inspect one bounty
./bin/au bounty show <id>

# move state
./bin/au bounty move <id> in_progress

# update fields
./bin/au bounty progress <id> 45
./bin/au bounty owner <id> alice
./bin/au bounty priority <id> high
./bin/au bounty next <id> "Open PR and attach proof"
./bin/au bounty block <id> "Waiting on wallet confirmation"

# append timestamped note to project output
./bin/au bounty note <id> "Posted social link, waiting for payout"

# grouped status board
./bin/au bounty board
```
