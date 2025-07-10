package telegram

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/danialmd81/my-subscribtion/telegram/helpers"

	"github.com/PuerkitoBio/goquery"
)

var (
	client      = &http.Client{}
	maxMessages = 100
	configs     = map[string]string{
		"ss":       "",
		"vmess":    "",
		"trojan":   "",
		"vless":    "",
		"hysteria": "",
		"other":    "",
	}
	ConfigFileIds = map[string]int32{
		"ss":       0,
		"vmess":    0,
		"trojan":   0,
		"vless":    0,
		"hysteria": 0,
		"other":    0,
	}
)

func Run() {
	fmt.Println("[INFO] Telegram collector running...")

	fileData, err := helpers.ReadFileContent("telegram/channels.txt")
	if err != nil {
		fmt.Println("[FATAL ERROR] ", err)
		return
	}

	lines := strings.SplitSeq(fileData, "\n")
	for line := range lines {
		url := strings.TrimSpace(line)
		if url == "" || strings.HasPrefix(url, "#") {
			continue
		}
		url = helpers.ChangeUrlToTelegramWebUrl(url)

		resp := HttpRequest(url)
		if resp != nil {
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err == nil {
				err = resp.Body.Close()
				if err == nil {
					fmt.Println(" ")
					fmt.Println("---------------------------------------")
					fmt.Println("[INFO] Crawling ", url)
					CrawlForV2ray(doc, url)
					fmt.Println("[INFO] Crawled ", url+"!")
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
		// from latest to oldest mode :
		linesArr := strings.Split(lines, "\n")
		linesArr = helpers.Reverse(linesArr)
		lines = strings.Join(linesArr, "\n")
		lines = strings.TrimSpace(lines)
		helpers.WriteToFile(lines, "telegram/"+proto+".txt")
	}
	fmt.Println("[INFO] All Done :D")
}

func AddConfigNames(config string, configtype string) string {
	configs := strings.Split(config, "\n")
	newConfigs := ""
	for _, extractedConfig := range configs {
		extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
		if extractedConfig == "" {
			continue
		}
		switch {
		case strings.HasPrefix(extractedConfig, "vmess://"):
			extractedConfig = EditVmessPs(extractedConfig, configtype, true)
			if extractedConfig != "" {
				ConfigFileIds["vmess"] += 1
				newConfigs += extractedConfig + "\n"
			}
		case strings.HasPrefix(extractedConfig, "ss://"):
			ConfigFileIds["ss"] += 1
			newConfigs += extractedConfig + " - " + strconv.Itoa(int(ConfigFileIds["ss"])) + "\n"
		case strings.HasPrefix(extractedConfig, "trojan://"):
			ConfigFileIds["trojan"] += 1
			newConfigs += extractedConfig + " - " + strconv.Itoa(int(ConfigFileIds["trojan"])) + "\n"
		case strings.HasPrefix(extractedConfig, "vless://"):
			ConfigFileIds["vless"] += 1
			newConfigs += extractedConfig + " - " + strconv.Itoa(int(ConfigFileIds["vless"])) + "\n"
		case strings.HasPrefix(extractedConfig, "hysteria://"), strings.HasPrefix(extractedConfig, "hy2://"):
			ConfigFileIds["hysteria"] += 1
			newConfigs += extractedConfig + " - " + strconv.Itoa(int(ConfigFileIds["hysteria"])) + "\n"
		default:
			ConfigFileIds["other"] += 1
			newConfigs += extractedConfig + " - " + strconv.Itoa(int(ConfigFileIds["other"])) + "\n"
		}
	}
	return newConfigs
}

func CrawlForV2ray(doc *goquery.Document, channelLink string) {
	// here we are updating our DOM to include the x messages
	// in our DOM and then extract the messages from that DOM
	messages := doc.Find(".tgme_widget_message_wrap").Length()
	fmt.Println("Fetched message:", messages)
	link, exist := doc.Find(".tgme_widget_message_wrap .js-widget_message").Last().Attr("data-post")

	if messages < maxMessages && exist {
		number := strings.Split(link, "/")[1]
		doc = GetMessages(maxMessages, doc, number, channelLink)
	}

	// get only messages that are inside code or pre tag and check for v2ray configs
	doc.Find("code,pre").Each(func(j int, s *goquery.Selection) {
		messageText, _ := s.Html()
		str := strings.ReplaceAll(messageText, "<br/>", "\n")
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
		messageText = doc.Text()
		line := strings.TrimSpace(messageText)
		lines := strings.SplitSeq(line, "\n")
		for data := range lines {
			extractedConfigs := strings.Split(ExtractConfig(data, []string{}), "\n")
			for _, extractedConfig := range extractedConfigs {
				extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
				if extractedConfig == "" {
					continue
				}
				switch {
				case strings.HasPrefix(extractedConfig, "vmess://"):
					extractedConfig = EditVmessPs(extractedConfig, "vmess", false)
					if extractedConfig != "" {
						configs["vmess"] += extractedConfig + "\n"
					}
				case strings.HasPrefix(extractedConfig, "ss://"):
					configs["ss"] += extractedConfig + "\n"
				case strings.HasPrefix(extractedConfig, "trojan://"):
					configs["trojan"] += extractedConfig + "\n"
				case strings.HasPrefix(extractedConfig, "vless://"):
					configs["vless"] += extractedConfig + "\n"
				case strings.HasPrefix(extractedConfig, "hysteria://"), strings.HasPrefix(extractedConfig, "hy2://"):
					configs["hysteria"] += extractedConfig + "\n"
				default:
					configs["other"] += extractedConfig + "\n"
				}
			}
		}
	})
}

func ExtractConfig(Txt string, Tempconfigs []string) string {
	line := strings.TrimSpace(Txt)
	if line == "" {
		return strings.Join(Tempconfigs, "\n")
	}
	switch {
	case strings.HasPrefix(line, "ss://"):
		Tempconfigs = append(Tempconfigs, "\n"+line)
	case strings.HasPrefix(line, "vmess://"):
		Tempconfigs = append(Tempconfigs, "\n"+line)
	case strings.HasPrefix(line, "trojan://"):
		Tempconfigs = append(Tempconfigs, "\n"+line)
	case strings.HasPrefix(line, "vless://"):
		Tempconfigs = append(Tempconfigs, "\n"+line)
	case strings.HasPrefix(line, "hysteria://"), strings.HasPrefix(line, "hy2://"):
		Tempconfigs = append(Tempconfigs, "\n"+line)
	default:
		if line != "" {
			Tempconfigs = append(Tempconfigs, "\n"+line)
		}
	}
	return strings.Join(Tempconfigs, "\n")
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
					data["ps"] = " - " + strconv.Itoa(int(ConfigFileIds[fileName])) + "\n"
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
