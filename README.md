# elma
elma is simple outline monitoring.

# configure

```[
  {
    "target": "http://example.com",
    "rate": "5s"
  },
  {
    "target": "https://example.co.jp",
    "rate": "5s",
    "healthy_status_code_list": [
      200, 301, 302
    ]
  }
]
```
