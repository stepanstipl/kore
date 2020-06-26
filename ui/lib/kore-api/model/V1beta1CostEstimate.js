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
import V1beta1CostEstimateElement from './V1beta1CostEstimateElement';

/**
 * The V1beta1CostEstimate model module.
 * @module model/V1beta1CostEstimate
 * @version 0.0.1
 */
class V1beta1CostEstimate {
    /**
     * Constructs a new <code>V1beta1CostEstimate</code>.
     * @alias module:model/V1beta1CostEstimate
     */
    constructor() { 
        
        V1beta1CostEstimate.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>V1beta1CostEstimate</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1beta1CostEstimate} obj Optional instance to populate.
     * @return {module:model/V1beta1CostEstimate} The populated <code>V1beta1CostEstimate</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1beta1CostEstimate();

            if (data.hasOwnProperty('costElements')) {
                obj['costElements'] = ApiClient.convertToType(data['costElements'], [V1beta1CostEstimateElement]);
            }
            if (data.hasOwnProperty('maxCost')) {
                obj['maxCost'] = ApiClient.convertToType(data['maxCost'], 'Number');
            }
            if (data.hasOwnProperty('minCost')) {
                obj['minCost'] = ApiClient.convertToType(data['minCost'], 'Number');
            }
            if (data.hasOwnProperty('preparedAt')) {
                obj['preparedAt'] = ApiClient.convertToType(data['preparedAt'], 'String');
            }
            if (data.hasOwnProperty('typicalCost')) {
                obj['typicalCost'] = ApiClient.convertToType(data['typicalCost'], 'Number');
            }
        }
        return obj;
    }

/**
     * @return {Array.<module:model/V1beta1CostEstimateElement>}
     */
    getCostElements() {
        return this.costElements;
    }

    /**
     * @param {Array.<module:model/V1beta1CostEstimateElement>} costElements
     */
    setCostElements(costElements) {
        this['costElements'] = costElements;
    }
/**
     * @return {Number}
     */
    getMaxCost() {
        return this.maxCost;
    }

    /**
     * @param {Number} maxCost
     */
    setMaxCost(maxCost) {
        this['maxCost'] = maxCost;
    }
/**
     * @return {Number}
     */
    getMinCost() {
        return this.minCost;
    }

    /**
     * @param {Number} minCost
     */
    setMinCost(minCost) {
        this['minCost'] = minCost;
    }
/**
     * @return {String}
     */
    getPreparedAt() {
        return this.preparedAt;
    }

    /**
     * @param {String} preparedAt
     */
    setPreparedAt(preparedAt) {
        this['preparedAt'] = preparedAt;
    }
/**
     * @return {Number}
     */
    getTypicalCost() {
        return this.typicalCost;
    }

    /**
     * @param {Number} typicalCost
     */
    setTypicalCost(typicalCost) {
        this['typicalCost'] = typicalCost;
    }

}

/**
 * @member {Array.<module:model/V1beta1CostEstimateElement>} costElements
 */
V1beta1CostEstimate.prototype['costElements'] = undefined;

/**
 * @member {Number} maxCost
 */
V1beta1CostEstimate.prototype['maxCost'] = undefined;

/**
 * @member {Number} minCost
 */
V1beta1CostEstimate.prototype['minCost'] = undefined;

/**
 * @member {String} preparedAt
 */
V1beta1CostEstimate.prototype['preparedAt'] = undefined;

/**
 * @member {Number} typicalCost
 */
V1beta1CostEstimate.prototype['typicalCost'] = undefined;






export default V1beta1CostEstimate;

