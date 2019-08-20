(function () {
    class GeoHeader extends HTMLElement {
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

            var today = new Date();
            var dd = String(today.getDate()).padStart(2, '0');
            var mm = String(today.getMonth() + 1).padStart(2, '0'); //January is 0!
            var yyyy = today.getFullYear();

            today = mm + '/' + dd + '/' + yyyy;

            var version = 'Not Provided'; // override if set by the SDO

            //  still need  <span> Distribution org, </span>  <span> Release Date, </span>
            this.attachShadow({ mode: 'open' });
            this.shadowRoot.innerHTML = `
                    <div style="overflow-wrap: break-word;width=100%">
                       
                   <h3> ${obj["csdco:expedition"]} </h3>,
                    <br>
                          ${obj.name},   
                       
                    </div>
                      `;
        }
    }
    window.customElements.define('geodex-header', GeoHeader);
})();