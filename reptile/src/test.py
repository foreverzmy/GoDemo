from urllib.parse import urlencode

import requests_html

if __name__ == '__main__':
    name = '剑来'
    session = requests_html.HTMLSession()
    data = {
        'searchtype': 'articlename',
        'action': 'login',
        'searchkey': name,
        'submit': '&#160;搜&#160;&#160;索&#160;'
    }
    data_str = urlencode(data, encoding='gb2312')
    print(data_str)
    headers = {
        'Content-Type': 'application/x-www-form-urlencoded',
    }
    cookies = {
        'jieqiVisitTime': 'jieqiArticlesearchTime%3D1522302706'
    }
    r = session.post('http://www.biquge.vip/modules/article/search.php',
                     data=data_str,
                     cookies=cookies,
                     headers=headers)
    r.encoding = 'gb2312'
    print(r.text)
