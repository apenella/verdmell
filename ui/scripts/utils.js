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

var Event = {
		_listeners: {},

		on: function (listener, callback) {
				this._listeners[listener] = callback;
		},

		off: function (listener) {
			delete this._listeners[listener];
		},

		notify: function (listener, data) {
			for (var listener in this._listeners){
				this._listeners[listener](data);
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

_.extend(Model.prototype, Event);


//
// View Object
//-----------------------------------------------------------
var View = function (options) {
	_.extend(this, options); 
	this.id = _.uniqueId('view');
};

_.extend(View.prototype, Event);

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