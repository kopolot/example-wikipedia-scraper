(() => {
    const elements = document.querySelectorAll('.mw-search-result')
    const searchParams = new URLSearchParams(window.location.search)
    const limit = +searchParams.get("limit")
    const page = (+searchParams.get('offset') / limit) + 1
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