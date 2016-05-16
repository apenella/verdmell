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
// Event object
//-----------------------------------------------------------

var Events = function Event(sender) {
		this._sender = sender;
		this._listeners = {};
};

Event.prototype = {
		on : function (listener, callback) {
				this._listeners[listener] = callback;
		},

		off : function (listener) {
			delete this._listeners[listener];
		},

		notify : function (args) {
				var index;

				for (var listener in this._listeners){
					this._listeners[listener](this._sender, args);
				}
		}
};

//
// Model Object
//-----------------------------------------------------------
var Model = function (attributes) {
	this.id = _.uniqueId('model');
	this.attributes = attributes || {};		
};

Model.prototype.get = function(attr) {
	return this.attributes[attr];
};

Model.prototype.set = function(attrs){
	if (_.isObject(attrs)) {
		_.extend(this.attributes, attrs);
		this.change(attrs);
	}
	return this;
};

Model.prototype.toJSON = function(options) {
	return _.clone(this.attributes);
};

Model.prototype.change = function(attrs){
	this.notify(this.id + 'update', attrs);
}; 

_.extend(Model.prototype, Events);


//
// View Object
//-----------------------------------------------------------
var View = function (options) {
	_.extend(this, options); 
	this.id = _.uniqueId('view');
};

_.extend(View.prototype, Events);

//
// Controller Object
//-----------------------------------------------------------
var Controller = function(options) {
	_.extend(this, options); 
	this.id = _.uniqueId('controller');

	var parts, selector, eventType;
	if(this.events){
		_.each(this.events, function(method, eventName){
			parts = eventName.split('.');
			selector = parts[0];
			eventType = parts[1];
			$(selector)['on' + eventType] = this[method];
		}.bind(this));
	}		
};