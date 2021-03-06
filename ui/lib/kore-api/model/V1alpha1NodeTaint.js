/**
 * Kore API
 * Kore API provides the frontend API (kore.appvia.io)
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
 * The V1alpha1NodeTaint model module.
 * @module model/V1alpha1NodeTaint
 * @version 0.0.1
 */
class V1alpha1NodeTaint {
    /**
     * Constructs a new <code>V1alpha1NodeTaint</code>.
     * @alias module:model/V1alpha1NodeTaint
     */
    constructor() { 
        
        V1alpha1NodeTaint.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>V1alpha1NodeTaint</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1alpha1NodeTaint} obj Optional instance to populate.
     * @return {module:model/V1alpha1NodeTaint} The populated <code>V1alpha1NodeTaint</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1alpha1NodeTaint();

            if (data.hasOwnProperty('effect')) {
                obj['effect'] = ApiClient.convertToType(data['effect'], 'String');
            }
            if (data.hasOwnProperty('key')) {
                obj['key'] = ApiClient.convertToType(data['key'], 'String');
            }
            if (data.hasOwnProperty('value')) {
                obj['value'] = ApiClient.convertToType(data['value'], 'String');
            }
        }
        return obj;
    }

/**
     * @return {String}
     */
    getEffect() {
        return this.effect;
    }

    /**
     * @param {String} effect
     */
    setEffect(effect) {
        this['effect'] = effect;
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
    getValue() {
        return this.value;
    }

    /**
     * @param {String} value
     */
    setValue(value) {
        this['value'] = value;
    }

}

/**
 * @member {String} effect
 */
V1alpha1NodeTaint.prototype['effect'] = undefined;

/**
 * @member {String} key
 */
V1alpha1NodeTaint.prototype['key'] = undefined;

/**
 * @member {String} value
 */
V1alpha1NodeTaint.prototype['value'] = undefined;






export default V1alpha1NodeTaint;

