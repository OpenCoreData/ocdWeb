import {
	html,
	render
} from './lit-html.js';

(function () {
    class SimpleGet extends HTMLElement {
        constructor() {
            super();

            const resID = this.getAttribute('res-id');

            // CAUTION DEV / DEMO HACK..  comment out in production!!!!!!
            var newstr = resID.replace(/opencoredata.org/i, '192.168.2.89:9900');
            // console.log(newstr); 

            // GET test
            function tj_providers(id) {
                return fetch(id, {
                    headers: { 'Content-Type': 'application/ld+json', },
                    })
                    .then(function (response) {
                        return response.json();
                    })
                    .then(function (myJson) {
                        //  console.log(id);
                        // console.log(JSON.stringify(myJson));
                        // return JSON.stringify(myJson);
                        return myJson;
                    });
            }

            // GET test call...
            tj_providers(newstr).then((feature) => {
                this.attachShadow({ mode: 'open' });


                this.shadowRoot.innerHTML = `
                <div style="overflow-wrap: break-word;width=100%">
                    Feature: ${feature["http://opencoredata.org/voc/csdco/v1/hole_ID"]} <br>
                     
                     PI(s): ${feature["http://opencoredata.org/voc/csdco/v1/pi"]}   
                    
                     
                </div>
                  `;


                // var count = Object.keys(providers).length;
                // const itemTemplates = [];
                // var i;
                // for (i = 0; i < count; i++) {
                //     // console.log(providers[i].name)
                //     itemTemplates.push(  `${providers[i].name}`);
                //     // console.log(itemTemplates)
                // }

                // var h =  `<div>${itemTemplates}</div>`;
                // this.shadowRoot.innerHTML = `${h}` ;
                // this.shadowRoot.appendChild(this.cloneNode(h));
            });
        }
    }
    window.customElements.define('geodex-getfeature', SimpleGet);
})();


