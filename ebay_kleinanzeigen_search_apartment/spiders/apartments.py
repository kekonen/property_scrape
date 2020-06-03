# -*- coding: utf-8 -*-
import scrapy


class ApartmentsSpider(scrapy.Spider):
    name = 'apartments'
    allowed_domains = ['ebay-kleinanzeigen.de']
    start_urls = ['http://ebay-kleinanzeigen.de/']

    def parse(self, response):
        pass
