# My Subscription

This repository provides a collection of subscription and configuration files for various proxy protocols and Telegram channels. It is organized to help users easily find and use the required configuration files for their needs.

## Repository Structure

- `src/` - Main source directory
  - `subs/` - Contains text files for different proxy protocols
  - `telegram/` - Contains Telegram channel lists and protocol files
  - `all/` - Contains merged and tested files for each protocol

## Text Files Table

| Protocol    | Subscription File                                                                                                   | Telegram File                                                                          | All (Merged & Tested) File                                                                 |
|-------------|---------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------|
| Hysteria    | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/subs/hysteria.txt)                      | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/telegram/hysteria.txt)         | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/all/hysteria.txt)                |
| Other       | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/subs/other.txt)                         | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/telegram/other.txt)            | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/all/other.txt)                   |
| Shadowsocks | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/subs/ss.txt)                            | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/telegram/ss.txt)               | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/all/ss.txt)                      |
| Trojan      | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/subs/trojan.txt)                        | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/telegram/trojan.txt)           | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/all/trojan.txt)                  |
| VLESS       | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/subs/vless.txt)                         | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/telegram/vless.txt)            | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/all/vless.txt)                   |
| VMess       | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/subs/vmess.txt)                         | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/telegram/vmess.txt)            | [raw](https://raw.githubusercontent.com/danialmd81/my-subscribtion/main/src/all/vmess.txt)                   |

## To Do

- [x] Refactor the Telegram `mixed.txt` file: Separate and categorize each protocol (e.g., VMess, VLESS, etc.) into their own dedicated files for better organization and clarity.
- [x] Improve categorization logic: Enhance the process for identifying and sorting protocols in Telegram files to ensure each protocol is clearly distinguished and easy to find.
- [ ] Implement actual connection test logic for each protocol

## Usage

- Browse the table above to find the configuration or subscription file you need.
- Click the raw link to view or download the file directly.
- Use these files in your proxy client or for Telegram channel management as needed.

## License

This repository is provided for educational and personal use. Please respect the terms of use for any third-party services referenced.
