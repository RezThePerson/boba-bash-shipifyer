run

```js
(() => {
  const el = document.querySelector('[data-page]');
  if (!el) return console.error('No element with [data-page] found!');

  const fullData = JSON.parse(el.dataset.page);
  console.log('Full Inertia Data:', fullData);
  
  copy(fullData);
  console.log('Copied full JSON to clipboard!');
  
  return fullData;
})();
```

on bash.hackclub.com/organize/34 to copy all page data
