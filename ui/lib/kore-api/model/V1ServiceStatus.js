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
 * The V1ServiceStatus model module.
 * @module model/V1ServiceStatus
 * @version 0.0.1
 */
class V1ServiceStatus {
    /**
     * Constructs a new <code>V1ServiceStatus</code>.
     * @alias module:model/V1ServiceStatus
     */
    constructor() { 
        
        V1ServiceStatus.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>V1ServiceStatus</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1ServiceStatus} obj Optional instance to populate.
     * @return {module:model/V1ServiceStatus} The populated <code>V1ServiceStatus</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1ServiceStatus();

            if (data.hasOwnProperty('components')) {
                obj['components'] = ApiClient.convertToType(data['components'], [V1Component]);
            }
            if (data.hasOwnProperty('configuration')) {
                obj['configuration'] = ApiClient.convertToType(data['configuration'], Object);
            }
            if (data.hasOwnProperty('message')) {
                obj['message'] = ApiClient.convertToType(data['message'], 'String');
            }
            if (data.hasOwnProperty('plan')) {
                obj['plan'] = ApiClient.convertToType(data['plan'], 'String');
            }
            if (data.hasOwnProperty('providerData')) {
                obj['providerData'] = ApiClient.convertToType(data['providerData'], Object);
            }
            if (data.hasOwnProperty('providerID')) {
                obj['providerID'] = ApiClient.convertToType(data['providerID'], 'String');
            }
            if (data.hasOwnProperty('serviceAccessEnabled')) {
                obj['serviceAccessEnabled'] = ApiClient.convertToType(data['serviceAccessEnabled'], 'Boolean');
            }
            if (data.hasOwnProperty('status')) {
                obj['status'] = ApiClient.convertToType(data['status'], 'String');
            }
        }
        return obj;
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
     * @return {Object}
     */
    getConfiguration() {
        return this.configuration;
    }

    /**
     * @param {Object} configuration
     */
    setConfiguration(configuration) {
        this['configuration'] = configuration;
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
    getPlan() {
        return this.plan;
    }

    /**
     * @param {String} plan
     */
    setPlan(plan) {
        this['plan'] = plan;
    }
/**
     * @return {Object}
     */
    getProviderData() {
        return this.providerData;
    }

    /**
     * @param {Object} providerData
     */
    setProviderData(providerData) {
        this['providerData'] = providerData;
    }
/**
     * @return {String}
     */
    getProviderID() {
        return this.providerID;
    }

    /**
     * @param {String} providerID
     */
    setProviderID(providerID) {
        this['providerID'] = providerID;
    }
/**
     * @return {Boolean}
     */
    getServiceAccessEnabled() {
        return this.serviceAccessEnabled;
    }

    /**
     * @param {Boolean} serviceAccessEnabled
     */
    setServiceAccessEnabled(serviceAccessEnabled) {
        this['serviceAccessEnabled'] = serviceAccessEnabled;
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
 * @member {Array.<module:model/V1Component>} components
 */
V1ServiceStatus.prototype['components'] = undefined;

/**
 * @member {Object} configuration
 */
V1ServiceStatus.prototype['configuration'] = undefined;

/**
 * @member {String} message
 */
V1ServiceStatus.prototype['message'] = undefined;

/**
 * @member {String} plan
 */
V1ServiceStatus.prototype['plan'] = undefined;

/**
 * @member {Object} providerData
 */
V1ServiceStatus.prototype['providerData'] = undefined;

/**
 * @member {String} providerID
 */
V1ServiceStatus.prototype['providerID'] = undefined;

/**
 * @member {Boolean} serviceAccessEnabled
 */
V1ServiceStatus.prototype['serviceAccessEnabled'] = undefined;

/**
 * @member {String} status
 */
V1ServiceStatus.prototype['status'] = undefined;






export default V1ServiceStatus;

