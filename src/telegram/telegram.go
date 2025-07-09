package telegram

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/danialmd81/my-subscribtion/telegram/helpers"

	"github.com/PuerkitoBio/goquery"
	"github.com/jszwec/csvutil"
)

var (
	client       = &http.Client{}
	maxMessages  = 100
	ConfigsNames = "@Vip_Security join us"
	configs      = map[string]string{
		"ss":       "",
		"vmess":    "",
		"trojan":   "",
		"vless":    "",
		"hysteria": "",
		"mixed":    "",
	}
	ConfigFileIds = map[string]int32{
		"ss":       0,
		"vmess":    0,
		"trojan":   0,
		"vless":    0,
		"hysteria": 0,
		"mixed":    0,
	}
	myregex = map[string]string{
		"ss":       `(?m)(...ss:|^ss:)\/\/.+?(%3A%40|#)`,
		"vmess":    `(?m)vmess:\/\/.+`,
		"trojan":   `(?m)trojan:\/\/.+?(%3A%40|#)`,
		"vless":    `(?m)vless:\/\/.+?(%3A%40|#)`,
		"hysteria": `(?m)(hysteria:\/\/|hy2:\/\/)[^\s]+`,
	}
	sort = flag.Bool("sort", false, "sort from latest to oldest (default : false)")
)

type ChannelsType struct {
	URL             string `csv:"URL"`
	AllMessagesFlag bool   `csv:"AllMessagesFlag"`
}

func Run() {

	flag.Parse()

	fileData, err := helpers.ReadFileContent("telegram/channels.csv")
	if err != nil {
		fmt.Println("[FATAL ERROR] ", err)
		return
	}
	var channels []ChannelsType
	if err = csvutil.Unmarshal([]byte(fileData), &channels); err != nil {
		fmt.Println("[FATAL ERROR] ", err)
		return
	}

	// loop through the channels lists
	for _, channel := range channels {

		// change url
		channel.URL = helpers.ChangeUrlToTelegramWebUrl(channel.URL)

		// get channel messages
		resp := HttpRequest(channel.URL)
		if resp != nil {
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err == nil {
				err = resp.Body.Close()
				if err == nil {
					fmt.Println(" ")
					fmt.Println("---------------------------------------")
					fmt.Println("[INFO] Crawling ", channel.URL)
					CrawlForV2ray(doc, channel.URL, channel.AllMessagesFlag)
					fmt.Println("[INFO] Crawled ", channel.URL+"!")
					fmt.Println("---------------------------------------")
					fmt.Println(" ")
				} else {
					fmt.Println("[ERROR] ", err)
				}
			} else {
				fmt.Println("[ERROR] Failed to parse document: ", err)
			}
		}
	}

	fmt.Println("[INFO] Creating output files !")

	for proto, configcontent := range configs {
		lines := helpers.RemoveDuplicate(configcontent)
		lines = AddConfigNames(lines, proto)
		if *sort {
			// 		from latest to oldest mode :
			linesArr := strings.Split(lines, "\n")
			linesArr = helpers.Reverse(linesArr)
			lines = strings.Join(linesArr, "\n")
		} else {
			// 		from oldest to latest mode :
			linesArr := strings.Split(lines, "\n")
			linesArr = helpers.Reverse(linesArr)
			linesArr = helpers.Reverse(linesArr)
			lines = strings.Join(linesArr, "\n")
		}
		lines = strings.TrimSpace(lines)
		helpers.WriteToFile(lines, "telegram/"+proto+".txt")

	}
	fmt.Println("[INFO] All Done :D")

}

func AddConfigNames(config string, configtype string) string {
	configs := strings.Split(config, "\n")
	newConfigs := ""
	for protoRegex, regexValue := range myregex {

		for _, extractedConfig := range configs {

			re := regexp.MustCompile(regexValue)
			matches := re.FindStringSubmatch(extractedConfig)
			if len(matches) > 0 {
				extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
				if extractedConfig != "" {
					switch protoRegex {
					case "vmess":
						extractedConfig = EditVmessPs(extractedConfig, configtype, true)
						if extractedConfig != "" {
							newConfigs += extractedConfig + "\n"
						}
					case "ss":
						Prefix := strings.Split(matches[0], "ss://")[0]
						if Prefix == "" {
							ConfigFileIds[configtype] += 1
							newConfigs += extractedConfig + ConfigsNames + " - " + strconv.Itoa(int(ConfigFileIds[configtype])) + "\n"
						}
					default:

						ConfigFileIds[configtype] += 1
						newConfigs += extractedConfig + ConfigsNames + " - " + strconv.Itoa(int(ConfigFileIds[configtype])) + "\n"
					}
				}
			}

		}
	}
	return newConfigs
}

func CrawlForV2ray(doc *goquery.Document, channelLink string, HasAllMessagesFlag bool) {
	// here we are updating our DOM to include the x messages
	// in our DOM and then extract the messages from that DOM
	messages := doc.Find(".tgme_widget_message_wrap").Length()
	fmt.Println("Fetched message:", messages)
	link, exist := doc.Find(".tgme_widget_message_wrap .js-widget_message").Last().Attr("data-post")

	if messages < maxMessages && exist {
		number := strings.Split(link, "/")[1]
		doc = GetMessages(maxMessages, doc, number, channelLink)
	}

	// extract v2ray based on message type and store configs at [configs] map
	if HasAllMessagesFlag {
		// get all messages and check for v2ray configs
		doc.Find(".tgme_widget_message_text").Each(func(j int, s *goquery.Selection) {
			// For each item found, get the band and title
			messageText, _ := s.Html()
			str := strings.Replace(messageText, "<br/>", "\n", -1)
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
			messageText = doc.Text()
			line := strings.TrimSpace(messageText)
			lines := strings.Split(line, "\n")
			for _, data := range lines {
				extractedConfigs := strings.Split(ExtractConfig(data, []string{}), "\n")
				for _, extractedConfig := range extractedConfigs {
					extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
					if extractedConfig != "" {

						// check if it is vmess or not
						re := regexp.MustCompile(myregex["vmess"])
						matches := re.FindStringSubmatch(extractedConfig)

						if len(matches) > 0 {
							extractedConfig = EditVmessPs(extractedConfig, "mixed", false)
							if line != "" {
								configs["mixed"] += extractedConfig + "\n"
							}
						} else {
							configs["mixed"] += extractedConfig + "\n"
						}

					}
				}
			}
		})
	} else {
		// get only messages that are inside code or pre tag and check for v2ray configs
		doc.Find("code,pre").Each(func(j int, s *goquery.Selection) {
			messageText, _ := s.Html()
			str := strings.ReplaceAll(messageText, "<br/>", "\n")
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
			messageText = doc.Text()
			line := strings.TrimSpace(messageText)
			lines := strings.Split(line, "\n")
			for _, data := range lines {
				extractedConfigs := strings.Split(ExtractConfig(data, []string{}), "\n")
				for protoRegex, regexValue := range myregex {

					for _, extractedConfig := range extractedConfigs {

						re := regexp.MustCompile(regexValue)
						matches := re.FindStringSubmatch(extractedConfig)
						if len(matches) > 0 {
							extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
							if extractedConfig != "" {
								switch protoRegex {
								case "vmess":
									extractedConfig = EditVmessPs(extractedConfig, protoRegex, false)
									if extractedConfig != "" {
										configs[protoRegex] += extractedConfig + "\n"
									}
								case "ss":
									Prefix := strings.Split(matches[0], "ss://")[0]
									if Prefix == "" {
										configs[protoRegex] += extractedConfig + "\n"
									}
								default:

									configs[protoRegex] += extractedConfig + "\n"
								}

							}
						}

					}

				}
			}

		})
	}
}

func ExtractConfig(Txt string, Tempconfigs []string) string {

	// filename can be "" or mixed
	for protoRegex, regexValue := range myregex {
		re := regexp.MustCompile(regexValue)
		matches := re.FindStringSubmatch(Txt)
		extractedConfig := ""
		if len(matches) > 0 {
			switch protoRegex {
			case "ss":
				Prefix := strings.Split(matches[0], "ss://")[0]
				if Prefix == "" {
					extractedConfig = "\n" + matches[0]
				} else if Prefix != "vle" { //  (Prefix != "vme" && Prefix != "") always true!
					d := strings.Split(matches[0], "ss://")
					extractedConfig = "\n" + "ss://" + d[1]
				}
			case "vmess":
				extractedConfig = "\n" + matches[0]
			default:
				extractedConfig = "\n" + matches[0]
			}

			Tempconfigs = append(Tempconfigs, extractedConfig)
			Txt = strings.ReplaceAll(Txt, matches[0], "")
			ExtractConfig(Txt, Tempconfigs)
		}
	}
	d := strings.Join(Tempconfigs, "\n")
	return d
}

func EditVmessPs(config string, fileName string, AddConfigName bool) string {
	// Decode the base64 string
	if config == "" {
		return ""
	}
	slice := strings.Split(config, "vmess://")
	if len(slice) > 0 {
		decodedBytes, err := base64.StdEncoding.DecodeString(slice[1])
		if err == nil {
			// Unmarshal JSON into a map
			var data map[string]interface{}
			err = json.Unmarshal(decodedBytes, &data)
			if err == nil {
				if AddConfigName {
					ConfigFileIds[fileName] += 1
					data["ps"] = ConfigsNames + " - " + strconv.Itoa(int(ConfigFileIds[fileName])) + "\n"
				} else {
					data["ps"] = ""
				}

				// marshal JSON into a map
				jsonData, _ := json.Marshal(data)
				// Encode JSON to base64
				base64Encoded := base64.StdEncoding.EncodeToString(jsonData)

				return "vmess://" + base64Encoded
			}
		}
	}

	return ""
}

func loadMore(link string) *goquery.Document {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println("[ERROR] Failed to create request: ", err)
		return nil
	}
	fmt.Println(link)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR] Request failed: ", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status: ", resp.StatusCode)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("[ERROR] Failed to parse document: ", err)
		return nil
	}
	return doc
}

func HttpRequest(url string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[ERROR] When requesting to: %s [ERROR] : %s\n", url, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[FATAL ERROR] ", err.Error())
	}
	return resp
}

func GetMessages(length int, doc *goquery.Document, number string, channel string) *goquery.Document {
	x := loadMore(channel + "?before=" + number)
	if x == nil {
		fmt.Println("[INFO] loadMore returned nil, returning current doc")
		return doc
	}

	html2, _ := x.Html()
	reader2 := strings.NewReader(html2)
	doc2, _ := goquery.NewDocumentFromReader(reader2)

	doc.Find("body").AppendSelection(doc2.Find("body").Children())

	newDoc := goquery.NewDocumentFromNode(doc.Selection.Nodes[0])
	messages := newDoc.Find(".js-widget_message_wrap").Length()

	if messages > length {
		return newDoc
	} else {
		num, _ := strconv.Atoi(number)
		n := num - 21
		if n > 0 {
			ns := strconv.Itoa(n)
			return GetMessages(length, newDoc, ns, channel)
		} else {
			return newDoc
		}
	}
}
