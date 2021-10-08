#!/bin/bash

set -u
set -e
set -o pipefail

CHARACTERS="
	https://foursouls.com/wp-content/uploads/2021/10/b2-blue_baby-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cain-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-eden-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-eve-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-isaac-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-judas-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-lazarus-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-lilith-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-maggy-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-samson-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_forgotten-768x1047.png
"

ETERNAL="
	https://foursouls.com/wp-content/uploads/2021/10/b2-blood_lust-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-book_of_belial-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-forever_alone-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-incubus-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-lazarus_rags-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-sleight_of_hand-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_bone-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_curse-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_d6-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-yum_heart-768x1047.png
"

TREASURE="
	https://foursouls.com/wp-content/uploads/2021/10/b2-baby_haunt-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-battery_bum-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-belly_button-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-blank_card-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-bobs_brain-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-book_of_sin-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-boomerang-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-box-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-breakfast-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-brimstone-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-bum_friend-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-bum-bo-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cambion_conception-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-champion_belt-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-chaos-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-chaos_card-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-charged_baby-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cheese_grater-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-compost-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-contract_from_below-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-crystal_ball-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-curse_of_the_tower-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dads_lost_coin-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-daddy_haunt-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dark_bum-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dead_bird-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-decoy-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dinner-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-diplopia-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-donation_machine-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dry_baby-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-edens_blessing-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-empty_vessel-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-eye_of_greed-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-fanny_pack-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-finger-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-flush-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-glass_cannon-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-goat_head-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-godhead-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-golden_razor_blade-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-greeds_gullet-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-guppys_collar-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-guppys_head-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-guppys_paw-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-host_hat-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-ipecac-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-jawbone-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-lucky_foot-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-meat-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-mini_mush-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-modeling_clay-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-moms_box-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-moms_bra-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-moms_coin_purse-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-moms_purse-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-moms_razor-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-moms_shovel-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-monster_manual-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-monstros_tooth-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-mr_boom-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-mystery_sack-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-no-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-pandoras_box-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-pay_to_play-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-placebo-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-polydactyly-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-portable_slot_machine-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-potato_peeler-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-razor_blade-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-remote_detonator-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-restock-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-sack_head-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-sack_of_pennies-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-sacred_heart-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-shadow-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-shiny_rock-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-smelter-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-spider_mod-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-spoon_bender-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-starter_deck-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-steamy_sale-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-suicide_king-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-synthoil-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-tarot_cloth-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-tech_x-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_battery-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_blue_map-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_chest-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_compass-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_d10-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_d100-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_d20-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_d4-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_dead_cat-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_habit-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_map-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_midas_touch-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_polaroid-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_poop-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_relic-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_shovel-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-theres_options-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-trinity_shield-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-two_of_clubs-768x1047.png
"

LOOT="
	https://foursouls.com/wp-content/uploads/2021/10/b2-a_dime-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-a_nickel-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-a_penny-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-blank_rune-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-bloody_penny-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-bomb-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-broken_ankh-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-butter_bean-3-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cains_eye-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-counterfeit_penny-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-curved_horn-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dagaz-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dice_shard-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-ehwaz-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-four_cents-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-gold_bomb-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-golden_horseshoe-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-guppys_hairball-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-i_the_magician-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-ii_the_high_priestess-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-iii_the_empress-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-iv_the_emperor-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-ix_the_hermit-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-lil_battery-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-lost_soul-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-mega_battery-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-o_the_fool-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-pills-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-pills-3-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-pills-2-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-purple_heart-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-soul_heart-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-swallowed_penny-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-three_cents-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-two_cents-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-v_the_hierophant-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-vi_the_lovers-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-vii_the_chariot-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-viii_justice-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-x_wheel_of_fortune-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xi_strength-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xii_the_hanged_man-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xiii_death-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xiv_temperance-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xiv_the_tower-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xix_the_sun-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xv_the_devil-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xvii_the_stars-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xviii_the_moon-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xx_judgement-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xxi_the_world-768x1047.png
"

MONSTERS="
	https://foursouls.com/wp-content/uploads/2021/10/b2-ambush-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-big_spider-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-black_bony-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-boom_fly-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-carrion_queen-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-chest-2-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-chest-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-chub-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-clotty-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cod_worm-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-conjoined_fatty-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-conquest-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-curse_of_amnesia-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-curse_of_greed-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-curse_of_loss-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-curse_of_pain-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-curse_of_the_blind-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cursed_chest-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cursed_fatty-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cursed_gaper-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cursed_horf-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cursed_keeper_head-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cursed_moms_hand-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-cursed_psy_horf-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-daddy_long_legs-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dank_globin-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dark_chest-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dark_chest-2-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dark_one-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-death-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-delirium-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-devil_deal-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dinga-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dip-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-dople-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-envy-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-evil_twin-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-famine-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-fat_bat-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-fatty-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-fly-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-gemini-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-gluttony-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-gold_chest-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-gold_chest-2-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-greed-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-greed_event-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-greedling-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-gurdy-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-gurdy_jr-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-hanger-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-holy_dinga-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-holy_dip-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-holy_keeper_head-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-holy_moms_eye-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-holy_squirt-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-hopper-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-horf-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-i_can_see_forever-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-keeper_head-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-larry_jr-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-leaper-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-leech-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-little_horn-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-lust-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-mask_of_infamy-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-mega_fatty-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-mega_troll_bomb-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-mom-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-moms_dead_hand-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-moms_eye-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-moms_hand-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-monstro-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-mulliboom-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-mulligan-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-pale_fatty-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-peep-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-pestilence-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-pin-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-pooter-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-portal-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-pride-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-psy_horf-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-rag_man-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-rage_creep-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-red_host-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-ring_of_flies-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-satan-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-scolex-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-secret_room-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-shop_upgrade-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-sloth-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-spider-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-squirt-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-stoney-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-swarm_of_flies-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_bloat-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_duke_of_flies-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_haunt-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-the_lamb-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-trite-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-troll_bombs-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-war-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-we_need_to_go_deeper-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-wizoob-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-wrath-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-xl_fLoor-768x1047.png
"

BONUS_SOULS="
	https://foursouls.com/wp-content/uploads/2021/10/b2-soul_of_gluttony-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-soul_of_greed-768x1047.png
	https://foursouls.com/wp-content/uploads/2021/10/b2-soul_of_guppy-768x1047.png
"

mkdir -pv base_game
pushd base_game
	if [[ `echo ${CHARACTERS} | wc -w` -ge 1 ]]; then
		mkdir -pv characters
		pushd characters
			for item in ${CHARACTERS}; do
				wget -nv ${item}
			done
		popd
	fi

	if [[ `echo ${ETERNAL} | wc -w` -ge 1 ]]; then
		mkdir -pv eternal
		pushd eternal
			for item in ${ETERNAL}; do
				wget -nv ${item}
			done
		popd
	fi

	if [[ `echo ${TREASURE} | wc -w` -ge 1 ]]; then
		mkdir -pv treasure
		pushd treasure
			for item in ${TREASURE}; do
				wget -nv ${item}
			done
		popd
	fi

	if [[ `echo ${LOOT} | wc -w` -ge 1 ]]; then
		mkdir -pv loot
		pushd loot
			for item in ${LOOT}; do
				wget -nv ${item}
			done
		popd
	fi

	if [[ `echo ${MONSTERS} | wc -w` -ge 1 ]]; then
		mkdir -pv monsters
		pushd monsters
			for item in ${MONSTERS}; do
				wget -nv ${item}
			done
		popd
	fi

	if [[ `echo ${BONUS_SOULS} | wc -w` -ge 1 ]]; then
		mkdir -pv bonus_souls
		pushd bonus_souls
			for item in ${BONUS_SOULS}; do
				wget -nv ${item}
			done
		popd
	fi
popd
