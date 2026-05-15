(() => {
    const tmpUrl = document.querySelector('#t-permalink > a')?.href ?? "";
    const searchParams = new URLSearchParams(tmpUrl.split('?')[1] ?? "")
    return {
        site_name: "wikipedia.pl",
        url: window.location.href,
        content: document.getElementById('mw-content-text').innerHTML,
        title: document.querySelector('h1 .mw-page-title-main').innerText,
        text_field_1: document.getElementById('footer-info-lastmod').innerText,
        text_field_2: document.querySelector('.infobox').innerHTML,
        text_field_3: document.querySelector('.grafika img')?.src ?? "",
        external_id: searchParams.get('oldid')
    }
})()