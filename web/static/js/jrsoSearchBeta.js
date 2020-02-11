/* jshint esversion: 6 */
import {
	html,
	render
} from './lit-html.js';

// https://stackoverflow.com/questions/36840396/fetch-gives-an-empty-response-body

// event listeners
document.querySelector('#q').addEventListener('keyup', function (e) {
	if (e.keyCode === 13) {
		searchActions();
	}
});

document.querySelector('#update').addEventListener('click', searchActions);

document.addEventListener('readystatechange', () => {
	if (document.readyState == 'complete') stateChangeSearch();
});

function pageLoadSearch() {
	console.log("=======================  Window load =======================");
}

// popstate for nav buttons
window.onpopstate = event => {
	stateChangeSearch();
};

function stateChangeSearch() {
	let params = (new URL(location)).searchParams;

	let q = params.get('q');
	let qdo = document.querySelector('#q');
	qdo.value = q;

	let n = params.get('n');
	let s = params.get('s');

	// TODO if we nav back to NULL..   we need to clear the result DOM
	// otherwise we still see the last none NULL search results
	if ( q != null) {
		blazefulltext(q, n, s);
	}
}

// rename to clickSearch  (needs updated elsewhere to do that)
function searchActions() {
	// let params = (new URL(location)).searchParams;
	let q = document.querySelector('#q').value
	let s = document.querySelector('#s').value
	let n = document.querySelector('#nn').value

	updateURL();
	blazefulltext(q, n, s);
	// updateNav();   // update the paging UI
}

function updateURL() {
	let qdo = document.querySelector('#q');
	let ndo = document.querySelector('#n');
	let sdo = document.querySelector('#s');
	let ido = document.querySelector('#i');

	let params = new URLSearchParams(location.search.slice(1));
	params.set('q', qdo.value);
	params.set('n', ndo.value);
	params.set('s', sdo.value);
	params.set('i', ido.value);


	//window.history.replaceState({}, '', location.pathname + '?' + params);
	const state = {
		q: qdo
	}
	window.history.pushState({}, '', location.pathname + '?' + params);
}

function blazefulltext(q, n, s) {

	(async () => {

		var url = new URL("http://triplestore.opencoredata.org/blazegraph/namespace/jrso/sparql"),
		// var url = new URL("http://triplestore.opencoredata.org/blazegraph/sparql"),

			// var url = new URL("http://192.168.2.89:8080/blazegraph/sparql"),
			// params = { query: "SELECT * { ?s ?p ?o  } LIMIT 11" }

			params = {query: ` prefix jrso: <http://opencoredata.org/voc/janus/v1/> \
	prefix schema: <http://schema.org/> \
	SELECT ?subj ?p ?score ?url  ?type  ?name ?relto ?addtype ?description \
	WHERE { \
   ?lit bds:search \"${q}\" . \
   ?lit bds:matchAllTerms "false" . \
   ?lit bds:relevance ?score . \
   ?subj ?p ?lit . \
   BIND (?subj as ?s) \ 
   ?s rdf:type schema:Dataset . \
   ?s jrso:hasLeg "166" .\
   OPTIONAL {?s schema:name ?name .} \
   OPTIONAL {?s schema:isRelatedTo ?relto .} \
   OPTIONAL {?s schema:additionalType ?addtype . } \
   OPTIONAL {?s schema:url ?url . } \
   OPTIONAL {?s schema:description ?description . } \
 } \
ORDER BY DESC(?score)  \
LIMIT 250 ` }

		Object.keys(params).forEach(key => url.searchParams.append(key, params[key]))

		const rawResponse = await fetch(url, {
			method: 'GET',
			// mode: 'no-cors', // no-cors, *cors, same-origin
			// cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
			credentials: 'omit', // include, *same-origin, omit
			headers: {
				'Accept': 'application/sparql-results+json',
				'Content-Type': 'application/json'
			} // ,
			// body: JSON.stringify({ query: 'SELECT * { ?s ?p ?o  } LIMIT 1', format: 'json' })
		});

		const content = await rawResponse.json();
		// console.log(content);

		const el = document.querySelector('#container2');
		const navel = document.querySelector('#container1');
		const s1 = document.querySelector('#side1');
		render(showresults(content), el);
		render(projresults(content), s1);

	})();
}

function getSafe(fn) {
	try {
		return fn();
	} catch (e) {
		return undefined;
	}
}

function truncate(n, useWordBoundary) {
	if (this.length <= n) { return this; }
	var subString = this.substr(0, n - 1);
	return (useWordBoundary
		? subString.substr(0, subString.lastIndexOf(' '))
		: subString) + "...";
};

// lit-html constant
//  SELECT ?subj  ?p ?score  ?type  ?name ?relto ?addtype ?url  ?description \
const showresults = (content) => {
	console.log("-----------------------------------------------")
	console.log(content)

	var barval = content.results.bindings
	var count = Object.keys(barval).length;
	const itemTemplates = [];

	itemTemplates.push(html`<p>Data Files</p>`);

	for (var i = 0; i < count; i++) {

		// if (getSafe(() => barval[i].type.value)) {

		// 	console.log(barval[i].type.value)

		// 	if (barval[i].type.value == "http://www.schema.org/DataSet") {

				itemTemplates.push(html`<div style="margin-top:30px">`);

				if (getSafe(() => barval[i].relto.value)) {
					itemTemplates.push(html`<div>See project:
					<a target="_blank" href="/id/do/${barval[i].relto.value}">${barval[i].relto.value}</a> </div>`);
				}

				if (getSafe(() => barval[i].name.value) && getSafe(() => barval[i].url.value)) {
					itemTemplates.push(html`<p>Digital Object: <a href="${barval[i].url.value}">${barval[i].name.value}</a> </p>`);
				}



				if (getSafe(() => barval[i].description.value)) {
					itemTemplates.push(html`<div> Description: ${barval[i].description.value} </div>`);
				}

				if (getSafe(() => barval[i].addtype.value)) {
					itemTemplates.push(html`<div> File type: ${barval[i].addtype.value} </div>`);
				}

				if (getSafe(() => barval[i].score.value)) {
					itemTemplates.push(html`<div> score: ${barval[i].score.value} </div>`);
				}

		// 	}
		// }

		itemTemplates.push(html`</div>`);
	}

	return html`
	<div style="margin-top:30px">
	   ${itemTemplates}
    </div>
	`;
};

// lit-html constant
//  SELECT ?subj  ?p ?score  ?type  ?name ?relto ?addtype ?url  ?description \
const projresults = (content) => {

	console.log("-top of Research Project---")

	var barval = content.results.bindings
	var count = Object.keys(barval).length;
	const itemTemplates = [];

	itemTemplates.push(html`<p>Related Projects</p>`);

	for (var i = 0; i < count; i++) {
		console.log("-in loop ---")


		if (getSafe(() => barval[i].type.value)) {
			if (barval[i].type.value == "http://schema.org/ResearchProject") {

				console.log("-Research Project---")

				itemTemplates.push(html`<div style="margin-top:30px">`);

				if (getSafe(() => barval[i].name.value) && getSafe(() => barval[i].url.value)) {
					itemTemplates.push(html`<p> <a href="${barval[i].url.value}">${barval[i].name.value}</a> </p>`);
				}

				if (getSafe(() => barval[i].relto.value)) {
					itemTemplates.push(html`<div> ${barval[i].relto.value} </div>`);
				}

				if (getSafe(() => barval[i].description.value)) {
					var s = barval[i].description.value
					itemTemplates.push(html`<div> Description: ${truncate.apply(s, [100, true])} </div>`);
				}


				if (getSafe(() => barval[i].addtype.value)) {
					itemTemplates.push(html`<div> ${barval[i].addtype.value} </div>`);
				}

				if (getSafe(() => barval[i].score.value)) {
					itemTemplates.push(html`<div> score: ${barval[i].score.value} </div>`);
				}

			}
		}

		itemTemplates.push(html`</div>`);
	}

	return html`
	<div style="margin-top:30px">
	   ${itemTemplates}
    </div>
	`;
};


// const OLDshowresults = (content) => {
// 	console.log("-----------------------------------------------")
// 	console.log(content);

// 	return html`<div style="text-align:center;margin-top:50px;position:relative">
// 	<br> Results:  ${content}</div>`;
// }

/*

async function postData(url = '', data = {}) {
  // Default options are marked with *
  const response = await fetch(url, {
	  method: 'POST', // *GET, POST, PUT, DELETE, etc.
	  mode: 'cors', // no-cors, *cors, same-origin
	  cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
	  credentials: 'same-origin', // include, *same-origin, omit
	  headers: {
		  'Content-Type': 'application/json'
		  // 'Content-Type': 'application/x-www-form-urlencoded',
	  },
	  redirect: 'follow', // manual, *follow, error
	  referrer: 'no-referrer', // no-referrer, *client
	  body: JSON.stringify(data) // body data type must match "Content-Type" header
  });
  return await response.json(); // parses JSON response into native JavaScript objects
}

try {
  const data = await postData('http://example.com/answer', { answer: 42 });
  console.log(JSON.stringify(data)); // JSON-string from `response.json()` call
} catch (error) {
  console.error(error);
}

*/
