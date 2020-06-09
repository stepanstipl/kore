/**
 * Appvia Kore API
 * Kore API provides the frontend API for the Appvia Kore (kore.appvia.io)
 *
 * The version of the OpenAPI document: 0.0.1
 * Contact: info@appvia.io
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */

import ApiClient from '../ApiClient';

/**
 * The V1OptionalSecretKeySelector model module.
 * @module model/V1OptionalSecretKeySelector
 * @version 0.0.1
 */
class V1OptionalSecretKeySelector {
    /**
     * Constructs a new <code>V1OptionalSecretKeySelector</code>.
     * @alias module:model/V1OptionalSecretKeySelector
     * @param key {String} 
     * @param name {String} 
     */
    constructor(key, name) { 
        
        V1OptionalSecretKeySelector.initialize(this, key, name);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, key, name) { 
        obj['key'] = key;
        obj['name'] = name;
    }

    /**
     * Constructs a <code>V1OptionalSecretKeySelector</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1OptionalSecretKeySelector} obj Optional instance to populate.
     * @return {module:model/V1OptionalSecretKeySelector} The populated <code>V1OptionalSecretKeySelector</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1OptionalSecretKeySelector();

            if (data.hasOwnProperty('key')) {
                obj['key'] = ApiClient.convertToType(data['key'], 'String');
            }
            if (data.hasOwnProperty('name')) {
                obj['name'] = ApiClient.convertToType(data['name'], 'String');
            }
            if (data.hasOwnProperty('namespace')) {
                obj['namespace'] = ApiClient.convertToType(data['namespace'], 'String');
            }
            if (data.hasOwnProperty('optional')) {
                obj['optional'] = ApiClient.convertToType(data['optional'], 'Boolean');
            }
        }
        return obj;
    }

/**
     * @return {String}
     */
    getKey() {
        return this.key;
    }

    /**
     * @param {String} key
     */
    setKey(key) {
        this['key'] = key;
    }
/**
     * @return {String}
     */
    getName() {
        return this.name;
    }

    /**
     * @param {String} name
     */
    setName(name) {
        this['name'] = name;
    }
/**
     * @return {String}
     */
    getNamespace() {
        return this.namespace;
    }

    /**
     * @param {String} namespace
     */
    setNamespace(namespace) {
        this['namespace'] = namespace;
    }
/**
     * @return {Boolean}
     */
    getOptional() {
        return this.optional;
    }

    /**
     * @param {Boolean} optional
     */
    setOptional(optional) {
        this['optional'] = optional;
    }

}

/**
 * @member {String} key
 */
V1OptionalSecretKeySelector.prototype['key'] = undefined;

/**
 * @member {String} name
 */
V1OptionalSecretKeySelector.prototype['name'] = undefined;

/**
 * @member {String} namespace
 */
V1OptionalSecretKeySelector.prototype['namespace'] = undefined;

/**
 * @member {Boolean} optional
 */
V1OptionalSecretKeySelector.prototype['optional'] = undefined;






export default V1OptionalSecretKeySelector;
