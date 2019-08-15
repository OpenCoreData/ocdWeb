(function () {
    class GeoKeywords extends HTMLElement {
        constructor() {
            super();

            // need to think about calling jsonld.js and using
            // it to parse the graph
            var obj;
            var inputs = document.getElementsByTagName('script');
            for (var i = 0; i < inputs.length; i++) {
                if (inputs[i].type.toLowerCase() == 'application/ld+json') {
                    obj = JSON.parse(inputs[i].innerHTML);
                }
            }

            this.attachShadow({ mode: 'open' });
            if (obj.keywords == undefined) {
                obj.keywords = "No Keywords Available"
            }
            this.shadowRoot.innerHTML = `
                    <div>
                       
                        ${obj.keywords}
                    </div>
                      `;
        }
    }
    window.customElements.define('geodex-keywords', GeoKeywords);
})();
