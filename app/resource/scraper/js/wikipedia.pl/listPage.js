(() => {
    const elements = document.querySelectorAll('.mw-search-result')
    const searchParams = new URLSearchParams(window.location.search)
    const limit = +searchParams.get("limit") || 100
    const page = Math.floor((+searchParams.get('offset') || 0) / limit) + 1
    const Pages = [];
    elements.forEach((v, i) => {
        const heading = v.querySelector('.mw-search-result-heading >  a')
        const Page = {
            URL: heading.href,
            externalID: heading.title,
            no: i + 1,
            page: page
        }
        Pages[i] = Page
    })
    return {
        pages: Pages
    }
})()