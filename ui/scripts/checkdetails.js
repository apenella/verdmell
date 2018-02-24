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
var checkDetailsModel = new Model('checkDetailsModel',{
	
	_check: [],

	set: function(data) {
		// console.log('checkDetailsModel::set');
		this._checks = _.where(data, {selected: true});

		// console.log('checkDetailsModel::set','checks selected',_.indexOf(data,_.where(data, {selected: true})));

		checkDetailsModel.set(this._checks);
	},

	//
	//
	// subscriptions
	observe: function(model, f) {
		checkDetailsModel.on(model.id, this.id, function(){f(model);}.bind(this));
	},

	// observe check for changes
	observeChecks: function(model) {
		// console.log('checkDetailsModel::observeChecks',model.attributes.getChecks().length);
		if (checksModel.attributes.getChecks().length > 0) {
			checkDetailsModel.attributes.set(model.attributes.getChecks());
		}
	}
});

//
// View Object
//-----------------------------------------------------------
var checkDetailsView = new View({
	parent: ".datacontainer",
	
	render:function(model) {
		// console.log('checkDetailsView::render',model);

		if ( model.length > 0 ) {
			$('.checkdetailsbg').remove();

			$('<div/>', {
				class: 'checkdetailsbg'
			}).appendTo('body');

			$('<div/>', {
				class: 'checkdetailsclose',
				text: 'Close'
			}).appendTo('.checkdetailsbg');
			
			$('<div/>', {
				class: 'checkdetails'
			}).appendTo('.checkdetailsbg');

			$('<div/>', {
				class: 'checkdetailstitle',
				text: model[0].checks.name
			}).appendTo('.checkdetails');

			$('<div/>', {
				class: 'checkdetailscheck'
			}).appendTo('.checkdetails');

			$('<div/>', {
				class: 'checkdetailssample'
			}).appendTo('.checkdetails');

			checkDetailsView.check('.checkdetailscheck',model[0].checks);
			checkDetailsView.sample('.checkdetailssample',model[0].samples);

			checkDetailsController.setEvents();
		}

	},

	check: function(parent,data) {
		console.log('checkDetailsView::check', parent, data);
		// Name
			$('<div/>',{
				class: 'checkdetailsheader',
				text: 'Check'
			}).appendTo(parent);
			$('<div/>',{
				class: 'checkdetailsinfo',
				text: data.name
			}).appendTo(parent);
			// Description
			$('<div/>',{
				class: 'checkdetailsheader',
				text: 'Description'
			}).appendTo(parent);
			$('<div/>',{
				class: 'checkdetailsinfo',
				text: data.description
			}).appendTo(parent);
			// Command
			$('<div/>',{
				class: 'checkdetailsheader',
				text: 'Command'
			}).appendTo(parent);
			$('<div/>',{
				class: 'checkdetailsinfo',
				text: data.command
			}).appendTo(parent);
	},

	sample: function(parent, data) {
		// Sample date
		$('<div/>',{
			class: 'checkdetailsheader',
			text: 'Sample time'
		}).appendTo(parent);
		$('<div/>',{
			class: 'checkdetailsinfo',
			text: data.Sample.sampletime
		}).appendTo(parent);
		// Exit value
		$('<div/>',{
			class: 'checkdetailsheader',
			text: 'Exit value'
		}).appendTo(parent);
		$('<div/>',{
			class: 'checkdetailsinfo',
			text: data.Sample.exit
		}).appendTo(parent);
		// Exit value
		$('<div/>',{
			class: 'checkdetailsheader',
			text: 'Elapsed time (ns)'
		}).appendTo(parent);
		$('<div/>',{
			class: 'checkdetailsinfo',
			text: data.Sample.elapsedtime
		}).appendTo(parent);
		// Output
		$('<div/>',{
			class: 'checkdetailsheader',
			text: 'Output'
		}).appendTo(parent);
		$('<div/>',{
			class: 'checkdetailsinfo',
			text: data.Sample.output
		}).appendTo(parent);
		// Timestamp
		$('<div/>',{
			class: 'checkdetailsheader',
			text: 'Timestamp'
		}).appendTo(parent);
		$('<div/>',{
			class: 'checkdetailsinfo',
			text: data.Sample.timestamp
		}).appendTo(parent);
	},

	close: function () {
		$('.checkdetailsbg').remove();
	},

	observe: function(model) {
		checkDetailsView.on(model.id, this.id, function(model){
			// console.log('checkDetailsView::observe', model);
			checkDetailsView.render(model);
		});
	}

});

//
// Controller Object
//-----------------------------------------------------------
var checkDetailsController = new Controller({

	// menu data management
	model: checkDetailsModel,
	// menu view management
	view: checkDetailsView,

	events: {
		".close::click": "close",
		".checkdetailsclose::click": "close"
	},

	initialize: function() {

		//
		// set listeners
		//
		// subscribe to checksModel
		this.model.attributes.observe(checksModel, this.model.attributes.observeChecks);
		// subscribe for changes to checksDetailsModel
		this.view.observe(this.model);
	},

		// set events to items
	setEvents: function() {
		var parts, selector, eventType;
		if(this.events){
			_.each(this.events, function(method, eventName){
				parts = eventName.split('::');
				// get item from events property
				selector = parts[0];
				// get method from events property
				eventType = parts[1];
				// hold the event to item
				$(selector).on(eventType, this[method]);
			}.bind(this));
		}
	},

	close: function() {
		console.log('checkDetailsController::close');
		checkDetailsView.close();
	}

});