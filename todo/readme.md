## ToDo App

After completing the main functionality, I switched the CLI to use Cobra. In doing this, I decided to switch up the flags to commands since most of them are actions and flags should be used for options, not actions. The commands are similarly named and have aliases:

- `-l` -> `list`/`l`
- `-q` -> `query`/`q`
- `-a` -> `add`/`a`
- `-u` -> `update`/`u`
- `-d` -> `delete`/`d`
- `-s` -> `status`/`s`

The `-f` (or `--file`) flag can be used to change the database file from any command. Since Cobra's flags do not support shorthand that is more than one ASCII character, I decided to switch from `-db`.

The `makefile` was updated to reflect this new command structure.