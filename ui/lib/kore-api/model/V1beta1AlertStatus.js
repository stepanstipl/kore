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
import V1beta1AlertRule from './V1beta1AlertRule';

/**
 * The V1beta1AlertStatus model module.
 * @module model/V1beta1AlertStatus
 * @version 0.0.1
 */
class V1beta1AlertStatus {
    /**
     * Constructs a new <code>V1beta1AlertStatus</code>.
     * @alias module:model/V1beta1AlertStatus
     */
    constructor() { 
        
        V1beta1AlertStatus.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>V1beta1AlertStatus</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1beta1AlertStatus} obj Optional instance to populate.
     * @return {module:model/V1beta1AlertStatus} The populated <code>V1beta1AlertStatus</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1beta1AlertStatus();

            if (data.hasOwnProperty('archivedAt')) {
                obj['archivedAt'] = ApiClient.convertToType(data['archivedAt'], 'String');
            }
            if (data.hasOwnProperty('detail')) {
                obj['detail'] = ApiClient.convertToType(data['detail'], 'String');
            }
            if (data.hasOwnProperty('rule')) {
                obj['rule'] = V1beta1AlertRule.constructFromObject(data['rule']);
            }
            if (data.hasOwnProperty('silencedUntil')) {
                obj['silencedUntil'] = ApiClient.convertToType(data['silencedUntil'], 'String');
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
    getArchivedAt() {
        return this.archivedAt;
    }

    /**
     * @param {String} archivedAt
     */
    setArchivedAt(archivedAt) {
        this['archivedAt'] = archivedAt;
    }
/**
     * @return {String}
     */
    getDetail() {
        return this.detail;
    }

    /**
     * @param {String} detail
     */
    setDetail(detail) {
        this['detail'] = detail;
    }
/**
     * @return {module:model/V1beta1AlertRule}
     */
    getRule() {
        return this.rule;
    }

    /**
     * @param {module:model/V1beta1AlertRule} rule
     */
    setRule(rule) {
        this['rule'] = rule;
    }
/**
     * @return {String}
     */
    getSilencedUntil() {
        return this.silencedUntil;
    }

    /**
     * @param {String} silencedUntil
     */
    setSilencedUntil(silencedUntil) {
        this['silencedUntil'] = silencedUntil;
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
 * @member {String} archivedAt
 */
V1beta1AlertStatus.prototype['archivedAt'] = undefined;

/**
 * @member {String} detail
 */
V1beta1AlertStatus.prototype['detail'] = undefined;

/**
 * @member {module:model/V1beta1AlertRule} rule
 */
V1beta1AlertStatus.prototype['rule'] = undefined;

/**
 * @member {String} silencedUntil
 */
V1beta1AlertStatus.prototype['silencedUntil'] = undefined;

/**
 * @member {String} status
 */
V1beta1AlertStatus.prototype['status'] = undefined;






export default V1beta1AlertStatus;
