package webdriver

import (
	"fmt"
	"log"
	"parser/internal/zip"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

/* Тут вводить хост, порт, логин, пароль от прокси
parser/internal/zip/modelszip.go
*/
// Добавление прокси
func AddProxy() string {
	plugin_file := "proxy_auth_plugin.zip"
	zip.CreateZipFile(plugin_file)
	log.Println("Создали .zip файл с прокси")
	return plugin_file
}

// Запуск вебдрайвера селениум + хром + прокси
func StartWebDriver() (selenium.WebDriver, *selenium.Service) {
	//plugin_file := AddProxy() // раскомментировать если нужны прокси
	//log.Println("Забираем информацию из proxy_auth_plugin.zip и включаем прокси")
	caps := selenium.Capabilities{
		"browserName": "chrome"}
	// Настройки для хрома
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--headless=new", // Включает режим Headless для браузера
			"--disable-gpu",
			"--disable-blink-features=AutomationControlled",
			"--start-maximized",
			"--start-fullscreen",
			"--no-proxy-server", // убрать если ставите прокси
			"user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
		},
	}
	//chromeCaps.AddExtension(plugin_file) //раскомментировать если нужны прокси
	caps.AddChrome(chromeCaps)

	// Установите путь к исполняемому файлу ChromeDriver
	chromeDriverPath := "../chrome/chromedriver"
	// Запуск ChromeDriver
	chromeDriverService, _ := selenium.NewChromeDriverService(chromeDriverPath, 9515)
	// Создание WebDriver
	wd, _ := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
	return wd, chromeDriverService

}

// Остановка приложения
func StopApp(wd selenium.WebDriver, chromeService *selenium.Service) {
	wd.Quit()
	chromeService.Stop()
}
