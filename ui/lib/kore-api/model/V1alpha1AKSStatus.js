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
import V1Component from './V1Component';

/**
 * The V1alpha1AKSStatus model module.
 * @module model/V1alpha1AKSStatus
 * @version 0.0.1
 */
class V1alpha1AKSStatus {
    /**
     * Constructs a new <code>V1alpha1AKSStatus</code>.
     * @alias module:model/V1alpha1AKSStatus
     */
    constructor() { 
        
        V1alpha1AKSStatus.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>V1alpha1AKSStatus</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1alpha1AKSStatus} obj Optional instance to populate.
     * @return {module:model/V1alpha1AKSStatus} The populated <code>V1alpha1AKSStatus</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1alpha1AKSStatus();

            if (data.hasOwnProperty('caCertificate')) {
                obj['caCertificate'] = ApiClient.convertToType(data['caCertificate'], 'String');
            }
            if (data.hasOwnProperty('components')) {
                obj['components'] = ApiClient.convertToType(data['components'], [V1Component]);
            }
            if (data.hasOwnProperty('endpoint')) {
                obj['endpoint'] = ApiClient.convertToType(data['endpoint'], 'String');
            }
            if (data.hasOwnProperty('message')) {
                obj['message'] = ApiClient.convertToType(data['message'], 'String');
            }
            if (data.hasOwnProperty('status')) {
                obj['status'] = ApiClient.convertToType(data['status'], 'String');
            }
        }
        return obj;
    }

/**
     * @return {String}
     */
    getCaCertificate() {
        return this.caCertificate;
    }

    /**
     * @param {String} caCertificate
     */
    setCaCertificate(caCertificate) {
        this['caCertificate'] = caCertificate;
    }
/**
     * @return {Array.<module:model/V1Component>}
     */
    getComponents() {
        return this.components;
    }

    /**
     * @param {Array.<module:model/V1Component>} components
     */
    setComponents(components) {
        this['components'] = components;
    }
/**
     * @return {String}
     */
    getEndpoint() {
        return this.endpoint;
    }

    /**
     * @param {String} endpoint
     */
    setEndpoint(endpoint) {
        this['endpoint'] = endpoint;
    }
/**
     * @return {String}
     */
    getMessage() {
        return this.message;
    }

    /**
     * @param {String} message
     */
    setMessage(message) {
        this['message'] = message;
    }
/**
     * @return {String}
     */
    getStatus() {
        return this.status;
    }

    /**
     * @param {String} status
     */
    setStatus(status) {
        this['status'] = status;
    }

}

/**
 * @member {String} caCertificate
 */
V1alpha1AKSStatus.prototype['caCertificate'] = undefined;

/**
 * @member {Array.<module:model/V1Component>} components
 */
V1alpha1AKSStatus.prototype['components'] = undefined;

/**
 * @member {String} endpoint
 */
V1alpha1AKSStatus.prototype['endpoint'] = undefined;

/**
 * @member {String} message
 */
V1alpha1AKSStatus.prototype['message'] = undefined;

/**
 * @member {String} status
 */
V1alpha1AKSStatus.prototype['status'] = undefined;






export default V1alpha1AKSStatus;
