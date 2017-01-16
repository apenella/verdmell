/*
	Package 'ui' 
	-server
	-handler
	-router
	-routes

	-html/
	-images/
	-pages/
	-scripts/
	-style/

*/

//
// Model Object
//-----------------------------------------------------------
var checksModel = new Model('checksModel',{

	_nodes: [],
	_checks: [],

	_itemTypeClass: {
		"nodes": "item",
		"services": "detailitem"
	},

	setNodes: function(data) {
		this._nodes = [],

		_.each(data, function(node){
			checksModel.attributes.addNode(node.name, node.URL);
		});
	},

	addNode: function(node, url) {
		this._nodes.push(checksModel.attributes.createNode(node, url));
	},

	getNode: function(node) {
		return _.where(checksModel.attributes._nodes, {node: node});
	},

	createNode: function (node, url) {
		return { node: node, url: url, selected: false };
	},

	addCheck: function(name, checks, samples) {
		this._checks.push(checksModel.attributes.createCheck(name, checks, samples));
	},

	createCheck: function(name, checks, samples) {
		return { name: name, checks: checks, samples: samples};
	},

	clearChecks: function() {
		// console.log('checksModel::clearChecks');
		this._checks = [];
	},

	setChecks: function() {
		// console.log('checksModel::setChecks');
		checksModel.set(this._checks);
	},

	setSelectedNode: function(node) {
		//console.log('checksModel::setSelectedNode',node);
		_.each(this._nodes, function(item){
			// set selected
			if ( item.node == node) {
				item.selected = true;
			} else {
				// if any other item is set as selected, unselect it
				if (item.selected) {
					item.selected = false;
				}
			}
		});
		// console.log('checksModel::setSelectedNode',this._nodes);
	},

	observe: function(model) {
		checksModel.on(model.id, this.id, function(model){
			if (detailsModel.attributes.getSelectedItem().length) {
				// console.log('checksModel::observe',detailsModel.attributes.getSelectedItem()[0].content.checks);
				// get node data related to detailitem clicked		
				node = checksModel.attributes.getNode(detailsModel.attributes.getSelectedItem()[0][checksModel.attributes._itemTypeClass[detailsModel.attributes.getSelectedItem()[0].type]]);
				//console.log('checksModel::observe',node[0]);
				if ( node.length ) {
					//checksModel.attributes.retrieveChecks(node[0].url);
					$.getJSON(node[0].url+MainController.nodeuri, function(data){
						// clear checks content
						checksModel.attributes.clearChecks();
						// console.log('checksModel::observe',data.checks);
						$.each(detailsModel.attributes.getSelectedItem()[0].content.checks, function(index, check){
							// console.log('checksModel::observe',data.checks.checks.checks[check]);
							// console.log('checksModel::observe',data.samples.Samples[check]);
							checksModel.attributes.addCheck(data.checks.checks.checks[check].name,data.checks.checks.checks[check],data.samples.Samples[check]);
						});
						checksModel.attributes.setChecks();
					});
				}
			}
		});
	}

});

//
// View Object
//-----------------------------------------------------------
var checksView = new View({
	// place to append the items
	parent: ".checks",

	initialize: function() {
		$('<div/>', {
			class: 'checkslist'
		}).appendTo(checksView.parent);
	},

	render: function(model) {
		$(checksView.parent).empty();

			$('<div/>', {
				class: 'checkstitle',
				id: 'title',
				text: 'CHECKS'
			}).appendTo(checksView.parent);

		_.each(model, function(item){
			$('<div/>', {
				class: 'check',
				id: item.name,
				status: item.samples.Sample.exit,
				text: item.name
			}).appendTo(checksView.parent);
		});
	},

	observe: function(model) {
		checksModel.on(model.id, this.id, function(model){
			// console.log('checksView::observe', model);
			checksView.render(model);
		});
	}

});

//
// Controller Object
//-----------------------------------------------------------
var checksController = new Controller({
	// menu data management
	model: checksModel,
	// menu view management
	view: checksView,

	initialize: function(data) {

		//
		// set listeners
		//
		// subscribe to detailsModel
		this.model.attributes.observe(detailsModel);
		// subscribe for changes to checksModel
		this.view.observe(this.model);

		this.view.initialize();
		// set nodes to model
		this.model.attributes.setNodes(data['nodes']);
	},

	update: function(data) {
		this.model.attributes.setNodes(data['nodes']);
	}

});
