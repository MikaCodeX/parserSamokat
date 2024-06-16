package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

type Product struct {
	Order []struct {
		AdressShop       string   `json:"Адрес доставки"`
		CatalogName      string   `json:"Категория"`
		NameProduct      []string `json:"Товар"`
		LinkProductImage []string `json:"Ссылка на картинку товара"`
		LinkProduct      []string `json:"Ссылка на товар"`
	}
}

// Получение информации о товарах на сайте
func GetDataProduct(wd selenium.WebDriver, adress string) Product {

	order := Product{}
	var urlProductImage, urlProduct, productNameAndPrice []string //sliceOrder
	var textUrlResult, sProduct, catalogName, sPriceProductWithDicount, sPriceProductWithoutDicount, sItog, sLinkForProduct string
	var elementsLinkTrue selenium.WebElement
	var catalogProduct, nameCatalogProduct, priceProductWithDiscount, priceProductWithoutDiscount []selenium.WebElement
	var chetProduct, nameChetProduct, chetPriceProductWithDiscount, chetPriceProductWithoutDiscount int
	chet := 1
	chetName := 1

	for i := 3; i <= 20; i++ {
		checkNameCatalogProduct := fmt.Sprintf("div.CategorySection_root__6Ai7Z:nth-child(%d) > div:nth-child(1) > span:nth-child(1)", i)
		elementsNameCatalogProduct, _ := wd.FindElements(selenium.ByCSSSelector, checkNameCatalogProduct)
		if elementsNameCatalogProduct == nil {
			log.Printf("На попытке %d ничего нет или не существует", i)
			break
		}
		nameCatalogProduct = append(nameCatalogProduct, elementsNameCatalogProduct...)
		if nameChetProduct < len(nameCatalogProduct) {
			catalogName, _ = nameCatalogProduct[nameChetProduct].Text()
		} else {
			break
		}
		nameChetProduct++
		chetName++

		for c := 1; ; c++ {
			// Переменные для парсинга по CSSSelector
			checkCatalogProduct := fmt.Sprintf("div.CategorySection_root__6Ai7Z:nth-child(%d) > div:nth-child(2) > a:nth-child(%d) > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > span:nth-child(1)", i, c)
			checkPriceProductWithDiscount := fmt.Sprintf("div.CategorySection_root__6Ai7Z:nth-child(%d) > div:nth-child(2) > a:nth-child(%d) > div:nth-child(1) > div:nth-child(2) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > span:nth-child(1) > span:nth-child(2)", i, c)
			checkPriceProductWithoutDiscount := fmt.Sprintf("div.CategorySection_root__6Ai7Z:nth-child(%d) > div:nth-child(2) > a:nth-child(%d) > div:nth-child(1) > div:nth-child(2) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > span:nth-child(1) > span:nth-child(1)", i, c)
			checkLinkForImageWithDiscount := fmt.Sprintf("div.CategorySection_root__6Ai7Z:nth-child(%d) > div:nth-child(2) > a:nth-child(%d) > div:nth-child(1) > div:nth-child(1) > img:nth-child(2)", i, c)
			checkLinkForImageWithoutDiscount := fmt.Sprintf("div.CategorySection_root__6Ai7Z:nth-child(%d) > div:nth-child(2) > a:nth-child(%d) > div:nth-child(1) > div:nth-child(1) > img:nth-child(1)", i, c)
			checkLinkForProduct := fmt.Sprintf("div.CategorySection_root__6Ai7Z:nth-child(%d) > div:nth-child(2) > a:nth-child(%d)", i, c)
			// Кладем в переменные данные которые вытащили по селектору
			elementsCatalogAndProduct, _ := wd.FindElements(selenium.ByCSSSelector, checkCatalogProduct)                          // Название товара и каталога
			elementsPriceProductWithDiscount, _ := wd.FindElements(selenium.ByCSSSelector, checkPriceProductWithDiscount)         // Цена со скидкой
			elementsPriceProductWithoutDiscount, _ := wd.FindElements(selenium.ByCSSSelector, checkPriceProductWithoutDiscount)   // Цена без скидки
			elementsLinkForImageWithDiscount, err1 := wd.FindElement(selenium.ByCSSSelector, checkLinkForImageWithDiscount)       // ссылка для товара со скидкой
			elementsLinkForImageWithoutDiscount, err2 := wd.FindElement(selenium.ByCSSSelector, checkLinkForImageWithoutDiscount) // ссылка для товара без скидки
			elementsLinkForProduct, _ := wd.FindElement(selenium.ByCSSSelector, checkLinkForProduct)

			if err1 != nil && err2 != nil {
				chet++
				break
			} else {
				if elementsLinkForImageWithDiscount == nil {
					elementsLinkTrue = elementsLinkForImageWithoutDiscount
					src, _ := elementsLinkTrue.GetAttribute("src")
					textUrlResult = src

				} else if elementsLinkForImageWithoutDiscount == nil {
					elementsLinkTrue = elementsLinkForImageWithDiscount
					src, _ := elementsLinkTrue.GetAttribute("src")
					textUrlResult = src

				}
			}

			if elementsCatalogAndProduct == nil {
				break

			}

			priceProductWithoutDiscount = append(priceProductWithoutDiscount, elementsPriceProductWithoutDiscount...)
			priceProductWithDiscount = append(priceProductWithDiscount, elementsPriceProductWithDiscount...)
			catalogProduct = append(catalogProduct, elementsCatalogAndProduct...)

			if chetPriceProductWithDiscount < len(priceProductWithDiscount) {
				sPriceProductWithDicount, _ = priceProductWithDiscount[chetPriceProductWithDiscount].Text()
				chetPriceProductWithDiscount++
			} else {
				sPriceProductWithDicount = "скидка отсутствует"
			}
			if chetProduct < len(catalogProduct) && chetPriceProductWithoutDiscount < len(priceProductWithoutDiscount) {
				sPriceProductWithoutDicount, _ = priceProductWithoutDiscount[chetPriceProductWithoutDiscount].Text()
				sProduct, _ = catalogProduct[chetProduct].Text()
				sLinkForProductDefault, _ := elementsLinkForProduct.GetAttribute("href")
				sLinkForProduct = sLinkForProductDefault
				if sPriceProductWithDicount == "скидка отсутствует" {
					sItog = fmt.Sprintf("Название: %s. Цена (со скидкой): %s, цена без скидки: %s", sProduct, sPriceProductWithDicount, sPriceProductWithoutDicount)
				} else {
					sItog = fmt.Sprintf("Название: %s. Цена (со скидкой): %s, цена без скидки: %s ₽", sProduct, sPriceProductWithDicount, sPriceProductWithoutDicount)
				}

				urlProductImage = append(urlProductImage, textUrlResult)
				urlProduct = append(urlProduct, sLinkForProduct)
			} else {
				break
			}
			productNameAndPrice = append(productNameAndPrice, sItog)
			chetProduct++

			chetPriceProductWithoutDiscount++
			chet++

		}

		order.Order = append(order.Order, struct { // Добавляем информацию в струткуру
			AdressShop       string   "json:\"Адрес доставки\""
			CatalogName      string   "json:\"Категория\""
			NameProduct      []string "json:\"Товар\""
			LinkProductImage []string "json:\"Ссылка на картинку товара\""
			LinkProduct      []string "json:\"Ссылка на товар\""
		}{
			AdressShop:       adress,
			CatalogName:      catalogName,
			NameProduct:      productNameAndPrice,
			LinkProductImage: urlProductImage,
			LinkProduct:      urlProduct,
		})
		productNameAndPrice = nil // Очистил слайс
		urlProductImage = nil     // Очистил слайс
		urlProduct = nil          // Очистил слайс

	}

	return order
}

// Вводим адрес на сайте
func EnterAdress(wd selenium.WebDriver) string {
	wd.Get("https://samokat.ru/?ysclid=lxhy7ci1wq933071069")
	time.Sleep(8 * time.Second)
	scanner := bufio.NewScanner(os.Stdin)
	var city, street, adress string
	fmt.Println("Введите город доставки: (для теста скопируйте и вставьте `Самара`)")
	scanner.Scan()
	city = scanner.Text()
	fmt.Println("Введите адрес доставки: (для теста cкопируйте и вставьте `Дзержинского, 29а`)")
	scanner.Scan()
	street = scanner.Text()
	log.Println("Вводим город на сайте")
	time.Sleep(5 * time.Second)
	changeCity, _ := wd.FindElement(selenium.ByCSSSelector, ".AddressConfirmBadge_buttons__Ou9hW > div:nth-child(2) > button:nth-child(1) > span:nth-child(1) > span:nth-child(1) > span:nth-child(1)")
	changeCity.Click()
	time.Sleep(5 * time.Second)
	findCity, _ := wd.FindElement(selenium.ByCSSSelector, "._textInput_1frhv_1")
	findCity.SendKeys(city)
	time.Sleep(2 * time.Second)
	enterCity, _ := wd.FindElement(selenium.ByCSSSelector, ".Suggest_suggestItems__wnlQV")
	enterCity.Click()
	time.Sleep(2 * time.Second)
	log.Println("Вводим улицу на сайте")
	time.Sleep(2 * time.Second)
	findStreet, _ := wd.FindElement(selenium.ByCSSSelector, "div.Suggest_root__KuclW:nth-child(1) > div:nth-child(1) > input:nth-child(1)")
	findStreet.SendKeys(street)
	time.Sleep(2 * time.Second)
	enterStreet, _ := wd.FindElement(selenium.ByCSSSelector, ".Suggest_suggestItems__wnlQV")
	enterStreet.Click()
	time.Sleep(2 * time.Second)
	enterButton, _ := wd.FindElement(selenium.ByCSSSelector, "._button--size_m_10nio_88 > button:nth-child(1)")
	enterButton.Click()
	time.Sleep(2 * time.Second)
	adress = fmt.Sprintf("г. %s, ул. %s", city, street)
	log.Println("Успешно ввели адрес доставки")
	return adress
}

// Открытие категорий и парсинг
func OpenCategoryParser(wd selenium.WebDriver) {
	order := Product{}
	// Открытие страницы
	adress := EnterAdress(wd)

	// Открытие категории

	for i := 3; i <= 33; i = i + 2 {

		if i > 3 { // Ограничения на открытия каталога (открываю 1 каталог) каждый последующий i + 2 (i = 5,7,9....)
			log.Println("Для парсинга большего кол-ва каталогов, поставьте i > 5,7,9,.....,33 в parser.go")
			break
		}
		time.Sleep(2 * time.Second)
		s, _ := wd.FindElement(selenium.ByCSSSelector, fmt.Sprintf("div.CatalogTreeSectionCard_details__5JoSN:nth-child(%d) > div:nth-child(1) > div:nth-child(2)", i))
		catalogSelenium, _ := wd.FindElement(selenium.ByCSSSelector, fmt.Sprintf("div.CatalogTreeSectionCard_details__5JoSN:nth-child(%d) > span:nth-child(2)", i))
		catalog, _ := catalogSelenium.Text()
		s.Click()
		log.Printf("Открываем каталог: %s", catalog)

		for j := 1; j <= 7; j++ {

			if j > 3 { // Ограничения на открытие категорий (открываю 3 вкладки)
				log.Println("Для парсинга большего кол-ва категорий в каталоге, поставьте j > 4,5,6,7 в parser.go")
				break
			}
			time.Sleep(2 * time.Second)
			k, _ := wd.FindElement(selenium.ByCSSSelector, fmt.Sprintf("div.CatalogTreeSectionCard_categories__4uYFm:nth-child(%d) > a:nth-child(%d) > span:nth-child(1)", i+1, j))

			if k == nil {
				break
			} else {
				k.Click()
				category, _ := k.Text()
				time.Sleep(1 * time.Second)

				log.Printf("Парсим данные в категории: %s", category)
				ordertest := GetDataProduct(wd, adress)
				log.Println("Складываем информацию в структуру")
				order.Order = append(order.Order, ordertest.Order...)
			}

		}

	}
	log.Println("Остановка парсера.....")
	WriteToTxt(order)

}

func WriteToTxt(product Product) {

	f, err := os.Create("../products/products.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	for _, p := range product.Order {
		fmt.Fprintf(f, "Адрес доставки: %s\n", p.AdressShop)
		fmt.Fprintf(f, "Категория: %s\n", p.CatalogName)
		fmt.Fprintf(f, "Товары:\n")
		for _, name := range p.NameProduct {
			fmt.Fprintf(f, "- %s\n", name)
		}
		fmt.Fprintf(f, "Ссылки на картинки товаров:\n")
		for _, url := range p.LinkProductImage {
			fmt.Fprintf(f, "- %s\n", url)
		}
		fmt.Fprintf(f, "Ссылки на товары:\n")
		for _, url := range p.LinkProduct {
			fmt.Fprintf(f, "- %s\n", url)
		}
		fmt.Fprintf(f, "\n")
	}

	log.Println("Создали файл products.txt в папке products")
}
