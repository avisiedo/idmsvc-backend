#!/usr/bin/env python3
"""Script to print out the last version available of kafka
"""

import requests
from html.parser import HTMLParser

DEFAULT_TIMEOUT=30

url = 'https://kafka.apache.org/downloads'
r = requests.get(url, allow_redirects=True, timeout=DEFAULT_TIMEOUT)

class KafkaHTMLDownloadParser(HTMLParser):
    """Download and parse the Kafka html page
       to extract the version available.
    """
    def __init__(self):
        HTMLParser.__init__(self)
        self.stack = []
        self.version = ''

    def get_version(self):
        return self.version

    def handle_starttag(self, tag, attrs):

        self.stack.append(tag)
        # if tag == 'h3':
        #     print(self.stack)
        # print("Encountered a start tag:", tag)

    def handle_endtag(self, tag):
        num_items = len(self.stack)
        if num_items == 0:
            return
        del self.stack[num_items-1]
        # print("Encountered an end tag :", tag)

    def handle_data(self, data):
        # print("Encountered some data  :", data)
        if self.version == '':
            if self.stack == ['html', 'head', 'link', 'link',
                              'meta', 'meta', 'body', 'div',
                              'div', 'div', 'h3']:
                self.version = data
        return

parser = KafkaHTMLDownloadParser()
parser.feed(r.text)

print(parser.get_version())
