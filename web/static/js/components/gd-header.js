(function () {
    class GeoHeader extends HTMLElement {
        constructor() {
            super();


            // TODO this is where starting to do this iwth the JSON-LD
            // code would help.  I could look for the parameters regardless of 
            // prefixed or expanded format

            // need to think about calling jsonld.js and using
            // it to parse the graph
            var obj;
            var inputs = document.getElementsByTagName('script');
            for (var i = 0; i < inputs.length; i++) {
                if (inputs[i].type.toLowerCase() == 'application/ld+json') {
                    obj = JSON.parse(inputs[i].innerHTML);
                }
            }

            //  still need  <span> Distribution org, </span>  <span> Release Date, </span>
            this.attachShadow({ mode: 'open' });
            this.shadowRoot.innerHTML = `
                    <div style="overflow-wrap: break-word;width=100%">
                   <h3>${obj.name}    (${obj["csdco:expedition"]}) </h3>,
                   
                    </div>
                      `;
        }
    }
    window.customElements.define('geodex-header', GeoHeader);
})();
