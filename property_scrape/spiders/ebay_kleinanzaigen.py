# -*- coding: utf-8 -*-
import scrapy


class EbayKleinanzaigenSpider(scrapy.Spider):
    name = 'ebay_kleinanzaigen'
    allowed_domains = ['ebay-kleinanzeigen.de']
    start_urls = ['http://ebay-kleinanzeigen.de/']

    def parse(self, response):
        pass
