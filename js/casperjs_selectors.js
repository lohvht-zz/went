// https://github.com/casperjs/casperjs/blob/a3cd29f6f9b1e6b8a632b55b6ef9e621086cb1c8/modules/clientutils.js

/**
 * This entire File outlines how CasperJS pulls stuff from XPath, would be nice to replicate in golang using cdp/chromedp
 * More implementation details for CasperJS can also be pulled!
 */

/**
 * Shortcut to build an XPath selector object.
 *
 * @param  String  expression  The XPath expression
 * @return Object
 * @see    http://casperjs.org/selectors.html
 */
function selectXPath(expression) {
  "use strict";
  return {
    type: 'xpath',
    path: expression,
    toString: function () {
      return this.type + ' selector: ' + this.path;
    }
  };
}
exports.selectXPath = selectXPath;

/**
 * Checks if an element matching the provided DOM CSS3/XPath selector exists in
 * current page DOM.
 *
 * @param  String  selector  A DOM CSS3/XPath selector
 * @return Boolean
 */
Casper.prototype.exists = function exists(selector) {
  "use strict";
  this.checkStarted();
  return this.callUtils("exists", selector);
};

/**
 * Invokes a client side utils object method within the remote page, with arguments.
 *
 * @param  {String}   method  Method name
 * @return {...args}          Arguments
 * @return {Mixed}
 * @throws {CasperError}      If invokation failed.
 */
Casper.prototype.callUtils = function callUtils(method) {
  "use strict";
  var args = [].slice.call(arguments, 1);
  var result = this.evaluate(function (method, args) {
    return __utils__.__call(method, args);
  }, method, args);
  if (utils.isObject(result) && result.__isCallError) {
    throw new CasperError(f("callUtils(%s) with args %s thrown an error: %s",
      method, args, result.message));
  }
  return result;
};

/**
 * Client Utils
 */

// public members
this.options = options || {};
this.options.scope = this.options.scope || document;

/**
 * Checks if a given DOM element exists in remote page.
 *
 * @param  String  selector  CSS3 selector
 * @return Boolean
 */
this.exists = function exists(selector) {
  try {
    return this.findAll(selector).length > 0;
  } catch (e) {
    return false;
  }
};

/**
 * Finds all DOM elements matching by the provided selector.
 *
 * @param  String | Object   selector  CSS3 selector (String only) or XPath object
 * @param  HTMLElement|null  scope     Element to search child elements within
 * @return Array|undefined
 */
this.findAll = function findAll(selector, scope) {
  scope = scope instanceof HTMLElement ? scope : scope && this.findOne(scope) || this.options.scope;
  try {
    var pSelector = this.processSelector(selector);
    if (pSelector.type === 'xpath') {
      return this.getElementsByXPath(pSelector.path, scope);
    } else {
      return Array.prototype.slice.call(scope.querySelectorAll(pSelector.path));
    }
  } catch (e) {
    this.log('findAll(): invalid selector provided "' + selector + '":' + e, "error");
  }
};

/**
   * Finds a DOM element by the provided selector.
   *
   * @param  String | Object   selector  CSS3 selector (String only) or XPath object
   * @param  HTMLElement|null  scope     Element to search child elements within
   * @return HTMLElement|undefined
   */
this.findOne = function findOne(selector, scope) {
  scope = scope instanceof HTMLElement ? scope : scope && this.findOne(scope) || this.options.scope;
  try {
    var pSelector = this.processSelector(selector);
    if (pSelector.type === 'xpath') {
      return this.getElementByXPath(pSelector.path, scope);
    } else {
      return scope.querySelector(pSelector.path);
    }
  } catch (e) {
    this.log('findOne(): invalid selector provided "' + selector + '":' + e, "error");
  }
};

var SUPPORTED_SELECTOR_TYPES = ['css', 'xpath'];
/**
 * Processes a selector input, either as a string or an object.
 *
 * If passed an object, if must be of the form:
 *
 *     selectorObject = {
 *         type: <'css' or 'xpath'>,
 *         path: <a string>
 *     }
 *
 * @param  String|Object  selector  The selector string or object
 *
 * @return an object containing 'type' and 'path' keys
 */
this.processSelector = function processSelector(selector) {
  var selectorObject = {
    toString: function toString() {
      return this.type + ' selector: ' + this.path;
    }
  };
  if (typeof selector === "string") {
    // defaults to CSS selector
    selectorObject.type = "css";
    selectorObject.path = selector;
    return selectorObject;
  } else if (typeof selector === "object") {
    // validation
    if (!selector.hasOwnProperty('type') || !selector.hasOwnProperty('path')) {
      throw new Error("Incomplete selector object");
    } else if (SUPPORTED_SELECTOR_TYPES.indexOf(selector.type) === -1) {
      throw new Error("Unsupported selector type: " + selector.type);
    }
    if (!selector.hasOwnProperty('toString')) {
      selector.toString = selectorObject.toString;
    }
    return selector;
  }
  throw new Error("Unsupported selector type: " + typeof selector);
};

/**
 * Retrieves a single DOM element matching a given XPath expression.
 *
 * @param  String            expression  The XPath expression
 * @param  HTMLElement|null  scope       Element to search child elements within
 * @return HTMLElement or null
 */
this.getElementByXPath = function getElementByXPath(expression, scope) {
  scope = scope || this.options.scope;
  var a = document.evaluate(expression, scope, this.xpathNamespaceResolver, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null);
  if (a.snapshotLength > 0) {
    return a.snapshotItem(0);
  }
};

/**
 * Retrieves all DOM elements matching a given XPath expression.
 *
 * @param  String            expression  The XPath expression
 * @param  HTMLElement|null  scope       Element to search child elements within
 * @return Array
 */
this.getElementsByXPath = function getElementsByXPath(expression, scope) {
  scope = scope || this.options.scope;
  var nodes = [];
  var a = document.evaluate(expression, scope, this.xpathNamespaceResolver, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null);
  for (var i = 0; i < a.snapshotLength; i++) {
    nodes.push(a.snapshotItem(i));
  }
  return nodes;
};

var XPATH_NAMESPACE = {
  svg: 'http://www.w3.org/2000/svg',
  mathml: 'http://www.w3.org/1998/Math/MathML'
};

/**
 * Build the xpath namespace resolver to evaluate on document
 *
 * @param String        prefix   The namespace prefix
 * @return the resolve namespace or null
 */
this.xpathNamespaceResolver = function xpathNamespaceResolver(prefix) {
  return XPATH_NAMESPACE[prefix] || null;
};

