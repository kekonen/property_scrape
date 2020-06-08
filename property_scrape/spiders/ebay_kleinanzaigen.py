# -*- coding: utf-8 -*-
import scrapy, random, re

def random_string(l):
    return ''.join(['abcdefghijklmnopqrstuvwxyz1234567890'[random.randrange(1,36)] for i in range(l)])
        

class EbayKleinanzaigenSpider(scrapy.Spider):
    name = 'ebay_kleinanzaigen'
    allowed_domains = ['ebay-kleinanzeigen.de']
    # start_urls = ['http://ebay-kleinanzeigen.de/']
    categories = {
        # 's-wohnung-kaufen': 196,
        # 's-auf-zeit-wg': 199,
        # ## 's-ferienwohnung-ferienhaus': 275,
        # 's-grundstuecke-garten': 207,
        # 's-haus-kaufen': 208,
        # 's-haus-mieten': 205,
        's-wohnung-mieten': 203,
    }
    category_counter = dict([[k,0] for k in categories.keys()])

    max_page = 1

    def start_requests(self):
        
        for c in self.categories:
            for p in range(1, self.max_page + 1):
                url = f'https://www.ebay-kleinanzeigen.de/{c}/seite:{p}/c{self.categories[c]}'
                yield scrapy.Request(url=url, callback=self.parse_listing, meta={'category': c}, headers={"user-agent": random_string(random.randrange(100,150))}) # 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36'
        

    def parse_listing(self, response):
        print('Listing...')
        for a in response.css('li.lazyload-item h2 a::attr(href)').getall():
            yield scrapy.Request(url=response.urljoin(a), callback=self.parse_page, meta={'category': response.meta['category']}, headers={"user-agent": random_string(random.randrange(100,150))})

    def parse_page(self, response):
        print('Item...\n\n\n')
        x = re.match(r'https\:\/\/www\.ebay\-kleinanzeigen.de\/[^\/]+\/(?P<name>[^\/]+)\/(?P<id>\d+)-(?P<cid>\d+)-(?P<hz>\d+)', response.url)
        name, _id, cid, hz = x[1], x[2], x[3], x[4]
        self.category_counter[response.meta['category']] += 1

        title = response.css('#viewad-title::text').get().strip()
        price_raw = response.css('#viewad-price::text').get().strip()
        locality_raw = response.css('#viewad-locality::text').get().strip()
        date_raw = response.css('#viewad-extra-info span::text').get().strip()
        counter_raw = response.css('#viewad-cntr span::text').get()
        if not counter_raw:
            print(response.url)
        print(title, price_raw, locality_raw, date_raw, counter_raw.strip() if counter_raw else None)

